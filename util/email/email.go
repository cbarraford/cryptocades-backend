package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

type Emailer struct {
	printOnly bool
	mg        mailgun.Mailgun
	templates *template.Template
}

func DefaultEmailer(projectRoot string) (Emailer, error) {
	return NewEmailer(
		os.Getenv("MAILGUN_DOMAIN"),
		os.Getenv("MAILGUN_API_KEY"),
		os.Getenv("MAILGUN_PUBLIC_KEY"),
		projectRoot,
		os.Getenv("ENVIRONMENT") != "production" && os.Getenv("ENVIRONMENT") != "staging",
	)
}

func NewEmailer(domain, apikey, publicKey, projectRoot string, printOnly bool) (Emailer, error) {
	tmpl := template.New("loader")
	tmpl, err := tmpl.ParseFiles(
		path.Join(projectRoot, "/util/email/templates/confirm.html"),
		path.Join(projectRoot, "/util/email/templates/confirm.txt"),
		path.Join(projectRoot, "/util/email/templates/password_reset.html"),
		path.Join(projectRoot, "/util/email/templates/password_reset.txt"),
		path.Join(projectRoot, "/util/email/templates/daily-dominance.html"),
		path.Join(projectRoot, "/util/email/templates/daily-dominance.txt"),
	)
	if err != nil {
		return Emailer{}, err
	}
	return Emailer{
		mg:        mailgun.NewMailgun(domain, apikey, publicKey),
		printOnly: printOnly,
		templates: tmpl,
	}, nil
}

func (e *Emailer) SendMessage(to, from, subject, body string) error {
	if e.printOnly {
		log.Printf("To:%s \nFrom:%s \nSubject: %s \nBody: %s",
			to, from, subject, body)
		return nil
	}
	message := e.mg.NewMessage(from, subject, body, to)
	_, _, err := e.mg.Send(message)
	return err
}

type EmailTemplate struct {
	Subject          string
	ConfirmURL       string
	ReferralURL      string
	PasswordResetURL string
}

func (e *Emailer) SendHTML(to, from, subject, templateName string, data interface{}) error {

	var err error
	var html string
	txt := bytes.Buffer{}
	buf := bytes.Buffer{}

	txtTmpl := e.templates.Lookup(fmt.Sprintf("%s.txt", templateName))
	txtTmpl.Execute(&txt, data)

	if e.printOnly {
		log.Printf("To:%s \nFrom:%s \nSubject: %s \nBody: %s",
			to, from, subject, txt.String())
		return nil
	}

	htmlTmpl := e.templates.Lookup(fmt.Sprintf("%s.html", templateName))
	htmlTmpl.Execute(&buf, data)
	html = buf.String()

	message := e.mg.NewMessage(from, subject, txt.String(), to)
	if html != "" {
		message.SetHtml(html)
	}
	_, _, err = e.mg.Send(message)
	return err
}
