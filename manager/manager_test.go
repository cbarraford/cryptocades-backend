package manager

import (
	"testing"

	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store/jackpot"
)

func TestPackage(t *testing.T) { TestingT(t) }

type ManagerSuite struct {
	store jackpot.Store
}

var _ = Suite(&ManagerSuite{})

type mockCreateJackpotStore struct {
	jackpot.Dummy
	created bool
	records []jackpot.Record
}

func (m *mockCreateJackpotStore) Create(record *jackpot.Record) error {
	m.created = true
	return nil
}

func (m *mockCreateJackpotStore) GetActiveJackpots() ([]jackpot.Record, error) {
	return m.records, nil
}

func (s *ManagerSuite) TestJackpotCreation(c *C) {
	store := &mockCreateJackpotStore{}
	c.Check(store.created, Equals, false)
	store.records = []jackpot.Record{
		{Jackpot: 200},
	}

	c.Assert(ManageJackpots(store), IsNil)
	c.Check(store.created, Equals, false)

	store.records = nil
	c.Assert(ManageJackpots(store), IsNil)
	c.Check(store.created, Equals, true)
}
