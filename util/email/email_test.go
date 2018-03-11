package email

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type EmailSuite struct{}

var _ = Suite(&EmailSuite{})

func (s *EmailSuite) TestSendMessage(c *C) {
	emailer, err := DefaultEmailer("../..")
	c.Assert(err, IsNil)
	c.Assert(emailer.SendMessage("cbarraford@gmail.com", "noreply@crytocades.com", "plain-text test", "body of the message"), IsNil)
}

func (s *EmailSuite) TestSendHTML(c *C) {
	emailer, err := DefaultEmailer("../..")
	c.Assert(err, IsNil)
	data := EmailTemplate{
		Subject:     "Confirm your cryptocades account",
		ConfirmURL:  "google.com",
		ReferralURL: "apple.com",
	}
	c.Assert(emailer.SendHTML("cbarraford@gmail.com", "noreply@crytocades.com", data.Subject, "confirm", data), IsNil)
}
