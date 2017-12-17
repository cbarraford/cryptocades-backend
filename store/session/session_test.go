package session

import (
	"fmt"
	"testing"

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

	record, err = s.store.GetByToken(record.Token)
	c.Check(record.CreatedTime.UnixNano() < record.ExpireTime.UnixNano(), Equals, true)
}

func (s *DBSuite) TestAuthenticate(c *C) {
	record := Record{
		UserId: 5,
	}

	err := s.store.Create(&record, 30)
	c.Assert(err, IsNil)

	var id int64
	id, err = s.store.Authenticate(record.Token)
	c.Assert(err, IsNil)
	c.Check(id, Equals, int64(5))

	err = s.store.Create(&record, -1)
	c.Assert(err, IsNil)
	_, err = s.store.Authenticate(record.Token)
	c.Assert(err, ErrorMatches, "Token expired.")

	_, err = s.store.Authenticate("bogus")
	c.Assert(err, NotNil)
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
