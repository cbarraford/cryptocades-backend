package email

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type EmailSuite struct{}

var _ = Suite(&EmailSuite{})

func (s *EmailSuite) TestSendMessage(c *C) {
	emailer := DefaultEmailer()
	c.Assert(emailer.SendMessage("to@bobby.com", "from@bobby.com", "subject", "body of the message"), IsNil)
}
