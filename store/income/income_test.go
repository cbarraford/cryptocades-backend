package income

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
	red := test.EphemeralRedisStore(c)
	s.store = store{sqlx: db, redis: red}
}

func (s *DBSuite) TearDownTest(c *C) {
	query := fmt.Sprintf("Truncate %s", table)
	_, err := s.store.sqlx.Exec(query)
	c.Assert(err, IsNil)

	_, err = s.store.redis.Do("FLUSHALL")
	c.Assert(err, IsNil)
}

func (s *DBSuite) TestCreateRequirements(c *C) {
	record := Record{
		SessionId: "abcfef",
	}
	c.Assert(s.store.Create(&record), NotNil)

	record = Record{
		UserId: 5,
	}
	c.Assert(s.store.Create(&record), NotNil)
}

func (s *DBSuite) TestCreate(c *C) {
	t := time.Now().Add(-1 * time.Hour)
	record1 := Record{
		UserId:        5,
		GameId:        4,
		SessionId:     "abcdef",
		Amount:        60,
		PartialAmount: 10,
		UpdatedTime:   t,
	}
	c.Assert(s.store.Create(&record1), IsNil)

	record2 := Record{
		UserId:        5,
		GameId:        4,
		SessionId:     "abcdef",
		Amount:        40,
		PartialAmount: 25,
	}
	c.Assert(s.store.Create(&record2), IsNil)

	r, err := s.store.Get(1)
	c.Assert(err, IsNil)
	c.Check(r.GameId, Equals, int64(4))
	c.Check(r.UserId, Equals, int64(5))
	c.Check(r.Amount, Equals, 101)
	c.Check(r.PartialAmount, Equals, 15)
	c.Check(r.SessionId, Equals, "abcdef")
	c.Check(r.UpdatedTime.Unix() > t.Unix(), Equals, true)

	record3 := Record{
		UserId:        9,
		GameId:        10,
		SessionId:     "bcdef",
		PartialAmount: 120,
	}
	c.Assert(s.store.Create(&record3), IsNil)

	r, err = s.store.Get(3)
	c.Assert(err, IsNil)
	c.Check(r.GameId, Equals, int64(10))
	c.Check(r.UserId, Equals, int64(9))
	c.Check(r.Amount, Equals, 6)
	c.Check(r.PartialAmount, Equals, 0)
	c.Check(r.SessionId, Equals, "bcdef")
}

func (s *DBSuite) TestListByUser(c *C) {
	record := Record{
		UserId:    6,
		GameId:    5,
		Amount:    40,
		SessionId: "dkuwfls",
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		UserId:    6,
		GameId:    8,
		SessionId: "bkdjfut",
		Amount:    30,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		UserId:    9,
		GameId:    14,
		SessionId: "sdlfkutjg",
		Amount:    90,
	}
	c.Assert(s.store.Create(&record), IsNil)

	records, err := s.store.ListByUser(6)
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 2)
	c.Check(records[0].GameId, Equals, int64(5))
	c.Check(records[0].UserId, Equals, int64(6))
	c.Check(records[0].Amount, Equals, 40)
	c.Check(records[1].GameId, Equals, int64(8))
	c.Check(records[1].UserId, Equals, int64(6))
	c.Check(records[1].Amount, Equals, 30)

	records, err = s.store.ListByUser(9)
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 1)
	c.Check(records[0].GameId, Equals, int64(14))
	c.Check(records[0].UserId, Equals, int64(9))
	c.Check(records[0].Amount, Equals, 90)

}

func (s *DBSuite) TestUserIncome(c *C) {
	record := Record{
		GameId:    4,
		SessionId: "sign up",
		UserId:    5,
		Amount:    60,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		GameId:    4,
		SessionId: "referral-2",
		UserId:    5,
		Amount:    50,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		GameId:    4,
		SessionId: "sign up",
		UserId:    6,
		Amount:    40,
	}
	c.Assert(s.store.Create(&record), IsNil)

	spent, err := s.store.UserIncome(5)
	c.Assert(err, IsNil)
	c.Check(spent, Equals, 110)

	spent, err = s.store.UserIncome(9999)
	c.Assert(err, IsNil)
	c.Check(spent, Equals, 0)

	var count int
	count, err = s.store.CountBonuses(5, "referral")
	c.Assert(err, IsNil)
	c.Check(count, Equals, 1)

	count, err = s.store.CountBonuses(5, "test")
	c.Assert(err, IsNil)
	c.Check(count, Equals, 0)
}

func (s *DBSuite) TestUserIncomeRank(c *C) {
	record := Record{
		GameId:    4,
		SessionId: "sign up",
		UserId:    5,
		Amount:    60,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		GameId:    4,
		SessionId: "referral-2",
		UserId:    6,
		Amount:    50,
	}
	c.Assert(s.store.Create(&record), IsNil)

	record = Record{
		GameId:    4,
		SessionId: "sign up",
		UserId:    7,
		Amount:    40,
	}
	c.Assert(s.store.Create(&record), IsNil)

	spent, err := s.store.UserIncomeRank(5)
	c.Assert(err, IsNil)
	c.Check(spent, Equals, 3)

	spent, err = s.store.UserIncomeRank(9999)
	c.Assert(err, IsNil)
	c.Check(spent, Equals, 0)
}
