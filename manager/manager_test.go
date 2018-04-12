package manager

import (
	"testing"

	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store/boost"
	"github.com/cbarraford/cryptocades-backend/store/entry"
	"github.com/cbarraford/cryptocades-backend/store/jackpot"
	"github.com/cbarraford/cryptocades-backend/store/matchup"
	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/cbarraford/cryptocades-backend/util/email"
)

func TestPackage(t *testing.T) { TestingT(t) }

type ManagerSuite struct {
	store jackpot.Store
}

var _ = Suite(&ManagerSuite{})

type mockUserStore struct {
	user.Dummy
}

func (m *mockUserStore) Get(id int64) (user.Record, error) {
	return user.Record{
		Id:      id,
		BTCAddr: "abcd",
	}, nil
}

type mockCreateJackpotStore struct {
	jackpot.Dummy
	created     bool
	updated     bool
	winnerId    int64
	btc_address string
	records     []jackpot.Record
	incomplete  []jackpot.Record
}

func (m *mockCreateJackpotStore) Create(record *jackpot.Record) error {
	m.created = true
	return nil
}

func (m *mockCreateJackpotStore) Update(record *jackpot.Record) error {
	m.updated = true
	m.winnerId = record.WinnerId
	m.btc_address = record.WinnerBTCAddr
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

	userStore := &mockUserStore{}

	c.Assert(ManageJackpots(store, entryStore, userStore), IsNil)
	c.Check(store.created, Equals, false)
	c.Check(store.updated, Equals, true)
	c.Check(store.winnerId > 0, Equals, true)
	c.Check(store.btc_address, Equals, "abcd")

	store.records = nil
	c.Assert(ManageJackpots(store, entryStore, userStore), IsNil)
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

type mockBoostStore struct {
	boost.Dummy
	created bool
	userId  int64
}

func (m *mockBoostStore) Create(record *boost.Record) error {
	m.created = true
	m.userId = record.UserId
	return nil
}

type mockMatchupsStore struct {
	matchup.Dummy
}

func (m *mockMatchupsStore) GetTopPerformers(match string, offset, top int) (records []matchup.Record, err error) {
	return []matchup.Record{
		{UserId: 5},
	}, nil
}

func (s *ManagerSuite) TestRewardPerformers(c *C) {
	matchupStore := &mockMatchupsStore{}
	boostStore := &mockBoostStore{}
	userStore := &mockUserStore{}

	em, err := email.DefaultEmailer("..")
	c.Assert(err, IsNil)

	c.Assert(RewardPerformers(1, matchupStore, boostStore, userStore, em), IsNil)
	c.Check(boostStore.created, Equals, false)
	c.Check(boostStore.userId, Equals, int64(0))

	c.Assert(RewardPerformers(0, matchupStore, boostStore, userStore, em), IsNil)
	c.Check(boostStore.created, Equals, true)
	c.Check(boostStore.userId, Equals, int64(5))
}
