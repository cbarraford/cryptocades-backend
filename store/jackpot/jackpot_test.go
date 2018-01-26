package jackpot

import (
	"fmt"
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/test"
)

func TestPackage(t *testing.T) { TestingT(t) }

type DBSuite struct {
	store store
}

var _ = Suite(&DBSuite{})

func (s *DBSuite) SetUpSuite(c *C) {
	db := test.EphemeralPostgresStore(c)
	s.store = store{sqlx: db}
}

func (s *DBSuite) TearDownTest(c *C) {
	query := fmt.Sprintf("Truncate %s", table)
	_, err := s.store.sqlx.Exec(query)
	c.Assert(err, IsNil)
}

func (s *DBSuite) TestCreateRequirements(c *C) {
	record := Record{}
	c.Assert(s.store.Create(&record), NotNil)
}

func (s *DBSuite) TestCreate(c *C) {
	endtime := time.Now().UTC().AddDate(0, 0, 1)
	record := Record{
		EndTime: endtime,
	}
	c.Assert(s.store.Create(&record), IsNil)

	var err error

	record, err = s.store.Get(record.Id)
	c.Assert(err, IsNil)
	c.Check(record.EndTime.Unix(), Equals, endtime.Unix())
}

func (s *DBSuite) TestGet(c *C) {
	record := Record{
		EndTime: time.Now().UTC().AddDate(0, 0, 1),
	}
	c.Assert(s.store.Create(&record), IsNil)

	_, err := s.store.Get(record.Id)
	c.Assert(err, IsNil)
}

func (s *DBSuite) TestUpdate(c *C) {
	record := Record{
		EndTime: time.Now().UTC().AddDate(0, 0, 1),
	}
	c.Assert(s.store.Create(&record), IsNil)

	r, err := s.store.Get(record.Id)
	c.Assert(err, IsNil)
	r.Jackpot = 101
	r.WinnerId = 45

	c.Assert(s.store.Update(&r), IsNil)

	r, err = s.store.Get(r.Id)
	c.Assert(err, IsNil)
	c.Check(r.Jackpot, Equals, 101)
	c.Check(r.WinnerId, Equals, int64(45))
}

func (s *DBSuite) TestList(c *C) {
	record := Record{
		Jackpot: 200,
		EndTime: time.Now().UTC().AddDate(0, 0, 1),
	}
	c.Assert(s.store.Create(&record), IsNil)

	records, err := s.store.List()
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 1)
	r := records[0]
	c.Check(r.Jackpot, Equals, 200)
}

func (s *DBSuite) TestGetActiveJackpots(c *C) {
	record := Record{
		Jackpot: 500,
		EndTime: time.Now().UTC().AddDate(0, 0, 1),
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		Jackpot: 300,
		EndTime: time.Now().UTC().AddDate(0, 0, -1),
	}
	c.Assert(s.store.Create(&record), IsNil)

	records, err := s.store.GetActiveJackpots()
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 1)
	c.Check(records[0].Jackpot, Equals, 500)
}
