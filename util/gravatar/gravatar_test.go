package gravatar

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type GravatarSuite struct{}

var _ = Suite(&GravatarSuite{})

func (s *GravatarSuite) TestAvatar(c *C) {
	expected := "https://www.gravatar.com/avatar/d9290cc27176c6fc74f4002f40fc9db8?d=%2Fimg%2Favatars%2F1.png&s=256"
	actual := Avatar("philipp@zoonman.com", 256)

	c.Check(actual, Equals, expected)

	expected = "https://www.gravatar.com/avatar/d8f4a1993546cc4b850cde3599e27aec?d=%2Fimg%2Favatars%2F5.png&s=100"
	actual = Avatar("not found", 100)
	c.Check(actual, Equals, expected)
}
