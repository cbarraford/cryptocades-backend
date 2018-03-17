package util

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type UtilSuite struct{}

var _ = Suite(&UtilSuite{})

func (s *UtilSuite) TestRandSeq(c *C) {
	r1 := RandSeq(14, Letters)
	r2 := RandSeq(14, LowerLetters)
	c.Check(r1, HasLen, 14)
	c.Check(r2, HasLen, 14)
	c.Check(r1, Not(Equals), r2)
}

func (s *UtilSuite) TestValidateEmail(c *C) {
	c.Check(ValidateEmail("test@test.com"), IsNil)
	c.Check(ValidateEmail("testtest.com"), NotNil)
	c.Check(ValidateEmail("test+test@test.com"), NotNil)
	c.Check(ValidateEmail(""), NotNil)
}

func (s *UtilSuite) TestValidateUsername(c *C) {
	c.Check(ValidateUsername("slkf73-92h9"), IsNil)
	c.Check(ValidateUsername("928371740934830"), IsNil)
	c.Check(ValidateUsername("Aallingh64"), IsNil)
	c.Check(ValidateUsername("-bad"), NotNil)
	c.Check(ValidateUsername("bad one"), NotNil)
	c.Check(ValidateUsername(""), NotNil)
	c.Check(ValidateUsername("111111111111111111111111111111111111111"), IsNil)
	c.Check(ValidateUsername("1111111111111111111111111111111111111111"), NotNil)
}
