package test

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type StoreSuite struct{}

var _ = Suite(&StoreSuite{})

func (s *StoreSuite) TestEphemeralPostgres(c *C) {
	EphemeralPostgresStore(c)
}
