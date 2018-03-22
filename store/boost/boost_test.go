package boost

import (
	"fmt"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store/income"
	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/cbarraford/cryptocades-backend/test"
)

func TestPackage(t *testing.T) { TestingT(t) }

type DBSuite struct {
	store   store
	users   user.Store
	user    user.Record
	user2   user.Record
	incomes income.Store
}

var _ = Suite(&DBSuite{})

func (s *DBSuite) SetUpSuite(c *C) {
	db := test.EphemeralPostgresStore(c)
	red := test.EphemeralRedisStore(c)
	s.store = store{sqlx: db}

	s.incomes = income.NewStore(db, red)
	s.users = user.NewStore(db)

	s.user = user.Record{
		Username: "testuser",
		Email:    "test@cryptocades.com",
		Password: "password",
	}
	c.Assert(s.users.Create(&s.user), IsNil)

	s.user2 = user.Record{
		Username: "testuser2",
		Email:    "test2@cryptocades.com",
		Password: "password",
	}
	c.Assert(s.users.Create(&s.user2), IsNil)
}

func (s *DBSuite) TearDownTest(c *C) {
	query := fmt.Sprintf("Truncate %s", table)
	_, err := s.store.sqlx.Exec(query)
	c.Assert(err, IsNil)
}

func (s *DBSuite) TestCreateRequirements(c *C) {
	record := Record{}
	c.Assert(s.store.Create(&record), ErrorMatches, "User id must not be blank")
}

func (s *DBSuite) TestCreate(c *C) {
	record := Record{
		UserId:     s.user.Id,
		IncomeId:   4,
		Multiplier: 10,
	}
	c.Assert(s.store.Create(&record), IsNil)

	r, err := s.store.Get(record.Id)
	c.Assert(err, IsNil)
	c.Check(r.UserId, Equals, s.user.Id)
	c.Check(r.IncomeId, Equals, int64(0)) // income id DOES NOT take
	c.Check(r.Multiplier, Equals, 2)      // multiplier DOES NOT take
}

func (s *DBSuite) TestListByUser(c *C) {
	record := Record{
		UserId: s.user.Id,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		UserId: s.user.Id,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		UserId: s.user2.Id,
	}
	c.Assert(s.store.Create(&record), IsNil)

	records, err := s.store.ListByUser(s.user.Id)
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 2)
	c.Check(records[0].UserId, Equals, s.user.Id)
	c.Check(records[1].UserId, Equals, s.user.Id)

	records, err = s.store.ListByUser(s.user2.Id)
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 1)
	c.Check(records[0].UserId, Equals, s.user2.Id)

	records, err = s.store.ListByUser(45)
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 0)
}

func (s *DBSuite) TestAssign(c *C) {
	inc := income.Record{
		UserId:    s.user.Id,
		GameId:    1,
		SessionId: "testing",
	}
	c.Assert(s.incomes.Create(&inc), IsNil)

	inc2 := income.Record{
		UserId:    s.user2.Id,
		GameId:    1,
		SessionId: "testing 88",
	}
	c.Assert(s.incomes.Create(&inc2), IsNil)

	record := Record{
		UserId: inc.UserId,
	}
	c.Assert(s.store.Create(&record), IsNil)

	c.Assert(s.store.Assign(record.Id, 2), ErrorMatches, "This boost and income session is not owned by the same user.")

	c.Assert(s.store.Assign(record.Id, 1), IsNil)
	c.Assert(s.store.Assign(record.Id, 2), ErrorMatches, "This boost is already assigned to a previous game session.")
	c.Assert(s.store.Assign(record.Id, 5555), NotNil) // bogus income id
}
