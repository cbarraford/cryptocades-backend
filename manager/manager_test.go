package manager

import (
	"testing"

	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store/entry"
	"github.com/cbarraford/cryptocades-backend/store/jackpot"
)

func TestPackage(t *testing.T) { TestingT(t) }

type ManagerSuite struct {
	store jackpot.Store
}

var _ = Suite(&ManagerSuite{})

type mockCreateJackpotStore struct {
	jackpot.Dummy
	created    bool
	updated    bool
	winnerId   int64
	records    []jackpot.Record
	incomplete []jackpot.Record
}

func (m *mockCreateJackpotStore) Create(record *jackpot.Record) error {
	m.created = true
	return nil
}

func (m *mockCreateJackpotStore) Update(record *jackpot.Record) error {
	m.updated = true
	m.winnerId = record.WinnerId
	return nil
}

func (m *mockCreateJackpotStore) GetActiveJackpots() ([]jackpot.Record, error) {
	return m.records, nil
}

func (m *mockCreateJackpotStore) GetIncompleteJackpots() ([]jackpot.Record, error) {
	return m.incomplete, nil
}

func (s *ManagerSuite) TestJackpotCreation(c *C) {
	store := &mockCreateJackpotStore{
		records: []jackpot.Record{
			{Jackpot: 200},
		},
		incomplete: []jackpot.Record{
			{Id: 101},
		},
	}
	c.Check(store.created, Equals, false)
	c.Check(store.updated, Equals, false)
	entryStore := &mockEntryStore{
		jackpots: []entry.Record{
			{JackpotId: 5, UserId: 10, Amount: 1},
			{JackpotId: 5, UserId: 11, Amount: 2},
			{JackpotId: 5, UserId: 12, Amount: 3},
			{JackpotId: 5, UserId: 13, Amount: 2},
			{JackpotId: 5, UserId: 14, Amount: 1},
		},
	}

	c.Assert(ManageJackpots(store, entryStore), IsNil)
	c.Check(store.created, Equals, false)
	c.Check(store.updated, Equals, true)
	c.Check(store.winnerId > 0, Equals, true)

	store.records = nil
	c.Assert(ManageJackpots(store, entryStore), IsNil)
	c.Check(store.created, Equals, true)
}

type mockEntryStore struct {
	entry.Dummy
	jackpots []entry.Record
}

func (m *mockEntryStore) ListByJackpot(id int64) ([]entry.Record, error) {
	return m.jackpots, nil
}

func (s *ManagerSuite) TestPickWinner(c *C) {
	store := &mockEntryStore{
		jackpots: []entry.Record{
			{JackpotId: 5, UserId: 10, Amount: 1},
			{JackpotId: 5, UserId: 11, Amount: 2},
			{JackpotId: 5, UserId: 12, Amount: 3},
			{JackpotId: 5, UserId: 13, Amount: 2},
			{JackpotId: 5, UserId: 14, Amount: 1},
		},
	}

	// run it a thousand times to check that we don't get some rare zero case
	for i := 0; i < 1000; i++ {
		winner, err := PickWinner(store, 5)
		c.Assert(err, IsNil)
		c.Check(winner > 0, Equals, true)
	}

	// test with no entries
	store.jackpots = nil
	winner, err := PickWinner(store, 5)
	c.Assert(err, IsNil)
	c.Check(winner > 0, Equals, false)
}

func (s *ManagerSuite) TestFindWinner(c *C) {
	store := &mockEntryStore{
		jackpots: []entry.Record{
			{JackpotId: 5, UserId: 10, Amount: 1},
			{JackpotId: 5, UserId: 11, Amount: 2},
			{JackpotId: 5, UserId: 12, Amount: 3},
			{JackpotId: 5, UserId: 13, Amount: 2},
			{JackpotId: 5, UserId: 14, Amount: 1},
		},
	}

	records, _ := store.ListByJackpot(1)
	c.Check(findWinner(records, 0).UserId, Equals, int64(0))
	c.Check(findWinner(records, 1).UserId, Equals, int64(10))
	c.Check(findWinner(records, 2).UserId, Equals, int64(11))
	c.Check(findWinner(records, 3).UserId, Equals, int64(11))
	c.Check(findWinner(records, 4).UserId, Equals, int64(12))
	c.Check(findWinner(records, 5).UserId, Equals, int64(12))
	c.Check(findWinner(records, 6).UserId, Equals, int64(12))
	c.Check(findWinner(records, 7).UserId, Equals, int64(13))
	c.Check(findWinner(records, 8).UserId, Equals, int64(13))
	c.Check(findWinner(records, 9).UserId, Equals, int64(14))
	c.Check(findWinner(records, 1000).UserId, Equals, int64(0))
}
