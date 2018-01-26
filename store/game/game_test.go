package game

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type GameSuite struct {
	store store
}

var _ = Suite(&GameSuite{})

func (s *GameSuite) SetUpSuite(c *C) {
	s.store = store{}
}

func (s *GameSuite) TestList(c *C) {
	games := s.store.List()
	c.Assert(games, HasLen, 2)
	c.Check(games[0].Id, Equals, 1)
	c.Check(games[0].Name, Equals, "Tallest Tower")
	c.Check(games[1].Id, Equals, 2)
	c.Check(games[1].Name, Equals, "Asteroid Tycoon")
}
