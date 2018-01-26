package entry

import (
	"fmt"
	"testing"

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
	record := Record{
		UserId: 5,
		Amount: 60,
	}
	c.Assert(s.store.Create(&record), NotNil)

	record = Record{
		JackpotId: 4,
		Amount:    60,
	}
	c.Assert(s.store.Create(&record), NotNil)

	record = Record{
		JackpotId: 4,
		UserId:    5,
	}
	c.Assert(s.store.Create(&record), NotNil)
}

func (s *DBSuite) TestCreate(c *C) {
	record1 := Record{
		JackpotId: 4,
		UserId:    5,
		Amount:    60,
	}
	c.Assert(s.store.Create(&record1), IsNil)

	record2 := Record{
		JackpotId: 4,
		UserId:    5,
		Amount:    40,
	}
	c.Assert(s.store.Create(&record2), IsNil)

	r, err := s.store.Get(record1.Id)
	c.Assert(err, IsNil)
	c.Check(r.JackpotId, Equals, int64(4))
	c.Check(r.UserId, Equals, int64(5))
	c.Check(r.Amount, Equals, 100)

	record3 := Record{
		JackpotId: 8,
		UserId:    9,
		Amount:    120,
	}
	c.Assert(s.store.Create(&record3), IsNil)

	r, err = s.store.Get(record3.Id)
	c.Assert(err, IsNil)
	c.Check(r.JackpotId, Equals, int64(8))
	c.Check(r.UserId, Equals, int64(9))
	c.Check(r.Amount, Equals, 120)
}

func (s *DBSuite) TestGetOdds(c *C) {
	record := Record{
		JackpotId: 4,
		UserId:    5,
		Amount:    60,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		JackpotId: 4,
		UserId:    6,
		Amount:    40,
	}
	c.Assert(s.store.Create(&record), IsNil)

	odd, err := s.store.GetOdds(record.JackpotId, record.UserId)
	c.Assert(err, IsNil)
	c.Check(odd.Total, Equals, int64(100))
	c.Check(odd.JackpotId, Equals, int64(4))
	c.Check(odd.Entries, Equals, int64(40))

	odd, err = s.store.GetOdds(500, 600)
	c.Assert(err, IsNil)
	c.Check(odd.Total, Equals, int64(0))
	c.Check(odd.JackpotId, Equals, int64(500))
	c.Check(odd.Entries, Equals, int64(0))
}

func (s *DBSuite) TestList(c *C) {
	record := Record{
		JackpotId: 4,
		UserId:    6,
		Amount:    40,
	}
	c.Assert(s.store.Create(&record), IsNil)

	records, err := s.store.List()
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 1)
	r := records[0]
	c.Check(r.JackpotId, Equals, int64(4))
	c.Check(r.UserId, Equals, int64(6))
	c.Check(r.Amount, Equals, 40)
}

func (s *DBSuite) TestListByUser(c *C) {
	record := Record{
		JackpotId: 4,
		UserId:    6,
		Amount:    40,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		JackpotId: 5,
		UserId:    6,
		Amount:    30,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		JackpotId: 6,
		UserId:    9,
		Amount:    90,
	}
	c.Assert(s.store.Create(&record), IsNil)

	records, err := s.store.ListByUser(6)
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 2)
	c.Check(records[0].JackpotId, Equals, int64(4))
	c.Check(records[0].UserId, Equals, int64(6))
	c.Check(records[0].Amount, Equals, 40)
	c.Check(records[1].JackpotId, Equals, int64(5))
	c.Check(records[1].UserId, Equals, int64(6))
	c.Check(records[1].Amount, Equals, 30)

	records, err = s.store.ListByUser(9)
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 1)
	c.Check(records[0].JackpotId, Equals, int64(6))
	c.Check(records[0].UserId, Equals, int64(9))
	c.Check(records[0].Amount, Equals, 90)

}

func (s *DBSuite) TestUserSpent(c *C) {
	record := Record{
		JackpotId: 4,
		UserId:    5,
		Amount:    60,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		JackpotId: 4,
		UserId:    5,
		Amount:    50,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		JackpotId: 4,
		UserId:    6,
		Amount:    40,
	}
	c.Assert(s.store.Create(&record), IsNil)

	spent, err := s.store.UserSpent(5)
	c.Assert(err, IsNil)
	c.Check(spent, Equals, 110)

	spent, err = s.store.UserSpent(9999)
	c.Assert(err, IsNil)
	c.Check(spent, Equals, 0)
}
