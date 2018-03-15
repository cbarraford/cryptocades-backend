package url

import (
	"log"
	"testing"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type URLSuite struct{}

var _ = Suite(&URLSuite{})

func (s *URLSuite) TestGet(c *C) {
	u := Get("/hello?referral=2309b8")
	log.Printf(u.String())
	c.Assert(u.String(), Equals, "http://localhost:3000/hello?referral=2309b8")
}
