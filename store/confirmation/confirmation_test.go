package confirmation

import (
	"fmt"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/CBarraford/lotto/test"

	"github.com/CBarraford/lotto/store/user"
)

func TestPackage(t *testing.T) { TestingT(t) }

type DBSuite struct {
	store     store
	userStore user.Store
}

var _ = Suite(&DBSuite{})

func (s *DBSuite) SetUpSuite(c *C) {
	db := test.EphemeralPostgresStore(c)
	red := test.EphemeralRedisStore(c)
	s.store = store{sqlx: db}
	s.userStore = user.NewStore(db, red)

	record := user.Record{
		Username: "bob",
		Email:    "bob@lotto.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.userStore.Create(&record), IsNil)
}

func (s *DBSuite) TearDownTest(c *C) {
	query := fmt.Sprintf("Truncate %s", table)
	_, err := s.store.sqlx.Exec(query)
	c.Assert(err, IsNil)
}

func (s *DBSuite) TestCreateRequirements(c *C) {
	record := Record{
		Email:  "bob@bob.com",
		UserId: 1,
	}
	c.Assert(s.store.Create(&record), NotNil)

	record = Record{
		Code:   "1233456",
		UserId: 1,
	}
	c.Assert(s.store.Create(&record), NotNil)

	record = Record{
		Code:  "1233456",
		Email: "bob@bob.com",
	}
	c.Assert(s.store.Create(&record), NotNil)
}

func (s *DBSuite) TestCreate(c *C) {
	record := Record{
		Code:   "123456",
		Email:  "bob@lotto.com",
		UserId: 1,
	}
	c.Assert(s.store.Create(&record), IsNil)
	c.Check(record.Code, Equals, "123456")
	c.Check(record.Email, Equals, "bob@lotto.com")
	c.Check(record.UserId, Equals, int64(1))
}

func (s *DBSuite) TestGet(c *C) {
	record := Record{
		Code:   "123456",
		Email:  "bob@lotto.com",
		UserId: 1,
	}
	c.Assert(s.store.Create(&record), IsNil)

	r, err := s.store.Get(record.Id)
	c.Assert(err, IsNil)
	c.Check(record.Code, Equals, r.Code)
	c.Check(record.Email, Equals, r.Email)
	c.Check(record.UserId, Equals, r.UserId)
}

func (s *DBSuite) TestGetByCode(c *C) {
	record := Record{
		Code:   "123456",
		Email:  "bob@lotto.com",
		UserId: 1,
	}
	c.Assert(s.store.Create(&record), IsNil)

	r, err := s.store.GetByCode(record.Code)
	c.Assert(err, IsNil)
	c.Check(record.Code, Equals, r.Code)
	c.Check(record.Email, Equals, r.Email)
	c.Check(record.UserId, Equals, r.UserId)
}

func (s *DBSuite) TestDelete(c *C) {
	var err error
	record := Record{
		Code:   "123456",
		Email:  "bob@lotto.com",
		UserId: 1,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record, err = s.store.Get(record.Id)
	c.Assert(err, IsNil)

	c.Assert(s.store.Delete(record.Id), IsNil)

	record, err = s.store.Get(record.Id)
	c.Assert(err, NotNil)

}
