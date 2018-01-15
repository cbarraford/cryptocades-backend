package session

import (
	"fmt"
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"github.com/CBarraford/lotto/test"
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

func (s *DBSuite) TestCreate(c *C) {
	record := Record{
		UserId: 5,
	}

	err := s.store.Create(&record, 30)
	c.Assert(err, IsNil)
	c.Check(record.Token, Not(Equals), "")
	c.Check(record.CreatedTime.Add(escalatedTime*time.Minute).UnixNano() == record.EscalatedTime.UnixNano(), Equals, true)

	record, err = s.store.GetByToken(record.Token)
	c.Check(record.CreatedTime.UnixNano() < record.ExpireTime.UnixNano(), Equals, true)
	c.Check(record.CreatedTime.Add(escalatedTime*time.Minute).UnixNano() == record.EscalatedTime.UnixNano(), Equals, true)
}

func (s *DBSuite) TestAuthenticate(c *C) {
	record := Record{
		UserId: 5,
	}

	err := s.store.Create(&record, 30)
	c.Assert(err, IsNil)

	var id int64
	var escalated bool
	id, escalated, err = s.store.Authenticate(record.Token)
	c.Assert(err, IsNil)
	c.Check(escalated, Equals, true)
	c.Check(id, Equals, int64(5))

	err = s.store.Create(&record, -1)
	c.Assert(err, IsNil)
	_, _, err = s.store.Authenticate(record.Token)
	c.Assert(err, ErrorMatches, "Token expired.")

	_, _, err = s.store.Authenticate("bogus")
	c.Assert(err, NotNil)

	record = Record{
		UserId:      5,
		CreatedTime: time.Now().UTC().Add(-30 * time.Minute),
	}

	err = s.store.Create(&record, 60)
	c.Assert(err, IsNil)

	id, escalated, err = s.store.Authenticate(record.Token)
	c.Assert(err, IsNil)
	c.Check(escalated, Equals, false)
	c.Check(id, Equals, int64(5))

}

func (s *DBSuite) TestDelete(c *C) {
	record := Record{
		UserId: 5,
	}

	err := s.store.Create(&record, 30)
	c.Assert(err, IsNil)

	err = s.store.Delete(record.Token)
	c.Assert(err, IsNil)

	_, err = s.store.GetByToken(record.Token)
	c.Assert(err, NotNil)
}
