package url

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type URLSuite struct{}

var _ = Suite(&URLSuite{})

func (s *URLSuite) TestGet(c *C) {
	u := Get("/hello")
	c.Assert(u.String(), Equals, "http://localhost:3000/hello")
}
