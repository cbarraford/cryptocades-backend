package email

import (
	"log"
	"os"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

type Emailer struct {
	printOnly bool
	mg        mailgun.Mailgun
}

func DefaultEmailer() *Emailer {
	return NewEmailer(
		os.Getenv("MAILGUN_DOMAIN"),
		os.Getenv("MAILGUN_API_KEY"),
		os.Getenv("MAILGUN_PUBLIC_KEY"),
		os.Getenv("ENVIRONMENT") != "production" && os.Getenv("ENVIRONMENT") != "staging",
	)
}

func NewEmailer(domain, apikey, publicKey string, printOnly bool) *Emailer {
	return &Emailer{
		mg:        mailgun.NewMailgun(domain, apikey, publicKey),
		printOnly: printOnly,
	}
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
