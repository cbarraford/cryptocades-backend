package gravatar

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type GravatarSuite struct{}

var _ = Suite(&GravatarSuite{})

func (s *GravatarSuite) TestAvatar(c *C) {
	expected := "https://www.gravatar.com/avatar/d9290cc27176c6fc74f4002f40fc9db8?s=256"
	actual := Avatar("philipp@zoonman.com", 256)

	c.Check(expected, Equals, actual)
}
