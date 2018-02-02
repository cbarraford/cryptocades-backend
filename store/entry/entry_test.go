package entry

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
	incomes income.Store
	users   user.Store
	user    user.Record
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

	in := income.Record{
		UserId:    s.user.Id,
		GameId:    1,
		Amount:    200,
		SessionId: "abcdef",
	}
	c.Assert(s.incomes.Create(&in), IsNil)
}

func (s *DBSuite) TearDownTest(c *C) {
	query := fmt.Sprintf("Truncate %s", table)
	_, err := s.store.sqlx.Exec(query)
	c.Assert(err, IsNil)
}

func (s *DBSuite) TestCreateRequirements(c *C) {
	record := Record{
		UserId: s.user.Id,
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
		UserId:    s.user.Id,
	}
	c.Assert(s.store.Create(&record), NotNil)
}

func (s *DBSuite) TestCreate(c *C) {
	record1 := Record{
		JackpotId: 4,
		UserId:    s.user.Id,
		Amount:    60,
	}
	c.Assert(s.store.Create(&record1), IsNil)

	record2 := Record{
		JackpotId: 4,
		UserId:    s.user.Id,
		Amount:    40,
	}
	c.Assert(s.store.Create(&record2), IsNil)

	r, err := s.store.Get(1)
	c.Assert(err, IsNil)
	c.Check(r.JackpotId, Equals, int64(4))
	c.Check(r.UserId, Equals, s.user.Id)
	c.Check(r.Amount, Equals, 100)

	user := user.Record{
		Username: "test5",
		Email:    "test5@test.com",
		Password: "password",
	}
	c.Assert(s.users.Create(&user), IsNil)

	in := income.Record{
		UserId:    user.Id,
		GameId:    1,
		Amount:    200,
		SessionId: "abcdef",
	}
	c.Assert(s.incomes.Create(&in), IsNil)

	record3 := Record{
		JackpotId: 8,
		UserId:    user.Id,
		Amount:    120,
	}
	c.Assert(s.store.Create(&record3), IsNil)

	r, err = s.store.Get(3)
	c.Assert(err, IsNil)
	c.Check(r.JackpotId, Equals, int64(8))
	c.Check(r.UserId, Equals, user.Id)
	c.Check(r.Amount, Equals, 120)
}

func (s *DBSuite) TestOverspend(c *C) {
	record1 := Record{
		JackpotId: 4,
		UserId:    s.user.Id,
		Amount:    100,
	}
	c.Assert(s.store.Create(&record1), IsNil)

	record2 := Record{
		JackpotId: 4,
		UserId:    s.user.Id,
		Amount:    140,
	}
	c.Assert(s.store.Create(&record2), ErrorMatches, "Insufficient funds.")
}

func (s *DBSuite) TestGetOdds(c *C) {
	record := Record{
		JackpotId: 4,
		UserId:    s.user.Id,
		Amount:    60,
	}
	c.Assert(s.store.Create(&record), IsNil)

	user := user.Record{
		Username: "test4",
		Email:    "test4@test.com",
		Password: "password",
	}
	c.Assert(s.users.Create(&user), IsNil)

	in := income.Record{
		UserId:    user.Id,
		GameId:    1,
		Amount:    200,
		SessionId: "abcdef",
	}
	c.Assert(s.incomes.Create(&in), IsNil)

	record = Record{
		JackpotId: 4,
		UserId:    user.Id,
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
		UserId:    s.user.Id,
		Amount:    40,
	}
	c.Assert(s.store.Create(&record), IsNil)

	records, err := s.store.List()
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 1)
	r := records[0]
	c.Check(r.JackpotId, Equals, int64(4))
	c.Check(r.UserId, Equals, s.user.Id)
	c.Check(r.Amount, Equals, 40)
}

func (s *DBSuite) TestListByUser(c *C) {
	record := Record{
		JackpotId: 4,
		UserId:    s.user.Id,
		Amount:    40,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		JackpotId: 5,
		UserId:    s.user.Id,
		Amount:    30,
	}
	c.Assert(s.store.Create(&record), IsNil)

	user2 := user.Record{
		Username: "test3",
		Email:    "test3@test.com",
		Password: "password",
	}
	c.Assert(s.users.Create(&user2), IsNil)

	in := income.Record{
		UserId:    user2.Id,
		GameId:    1,
		Amount:    200,
		SessionId: "abcdef",
	}
	c.Assert(s.incomes.Create(&in), IsNil)

	record = Record{
		JackpotId: 6,
		UserId:    user2.Id,
		Amount:    90,
	}
	c.Assert(s.store.Create(&record), IsNil)

	records, err := s.store.ListByUser(s.user.Id)
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 2)
	c.Check(records[0].JackpotId, Equals, int64(4))
	c.Check(records[0].UserId, Equals, s.user.Id)
	c.Check(records[0].Amount, Equals, 40)
	c.Check(records[1].JackpotId, Equals, int64(5))
	c.Check(records[1].UserId, Equals, s.user.Id)
	c.Check(records[1].Amount, Equals, 30)

	records, err = s.store.ListByUser(user2.Id)
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 1)
	c.Check(records[0].JackpotId, Equals, int64(6))
	c.Check(records[0].UserId, Equals, user2.Id)
	c.Check(records[0].Amount, Equals, 90)

}

func (s *DBSuite) TestListByJackpot(c *C) {
	record := Record{
		JackpotId: 4,
		UserId:    s.user.Id,
		Amount:    40,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		JackpotId: 5,
		UserId:    s.user.Id,
		Amount:    30,
	}
	c.Assert(s.store.Create(&record), IsNil)

	user2 := user.Record{
		Username: "test2",
		Email:    "test2@test.com",
		Password: "password",
	}
	c.Assert(s.users.Create(&user2), IsNil)

	in := income.Record{
		UserId:    user2.Id,
		GameId:    1,
		Amount:    200,
		SessionId: "abcdef",
	}
	c.Assert(s.incomes.Create(&in), IsNil)

	record = Record{
		JackpotId: 4,
		UserId:    user2.Id,
		Amount:    90,
	}
	c.Assert(s.store.Create(&record), IsNil)

	records, err := s.store.ListByJackpot(4)
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 2)
	c.Check(records[0].JackpotId, Equals, int64(4))
	c.Check(records[0].UserId, Equals, s.user.Id)
	c.Check(records[0].Amount, Equals, 40)
	c.Check(records[1].JackpotId, Equals, int64(4))
	c.Check(records[1].UserId, Equals, user2.Id)
	c.Check(records[1].Amount, Equals, 90)
}

func (s *DBSuite) TestUserSpent(c *C) {
	record := Record{
		JackpotId: 4,
		UserId:    s.user.Id,
		Amount:    60,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		JackpotId: 4,
		UserId:    s.user.Id,
		Amount:    50,
	}
	c.Assert(s.store.Create(&record), IsNil)

	user := user.Record{
		Username: "test7",
		Email:    "test7@test.com",
		Password: "password",
	}
	c.Assert(s.users.Create(&user), IsNil)

	in := income.Record{
		UserId:    user.Id,
		GameId:    1,
		Amount:    200,
		SessionId: "abcdef",
	}
	c.Assert(s.incomes.Create(&in), IsNil)

	record = Record{
		JackpotId: 4,
		UserId:    user.Id,
		Amount:    40,
	}
	c.Assert(s.store.Create(&record), IsNil)

	spent, err := s.store.UserSpent(s.user.Id)
	c.Assert(err, IsNil)
	c.Check(spent, Equals, 110)

	spent, err = s.store.UserSpent(9999)
	c.Assert(err, IsNil)
	c.Check(spent, Equals, 0)
}
