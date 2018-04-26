package matchup

import (
	"fmt"
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/cbarraford/cryptocades-backend/test"
	"github.com/garyburd/redigo/redis"
)

func TestPackage(t *testing.T) { TestingT(t) }

type DBSuite struct {
	store store
	users user.Store
	user  user.Record
	user2 user.Record
	user3 user.Record
}

var _ = Suite(&DBSuite{})

func (s *DBSuite) SetUpSuite(c *C) {
	db := test.EphemeralPostgresStore(c)
	red := test.EphemeralRedisStore(c)
	s.store = store{sqlx: db, redis: red}

	s.users = user.NewStore(db)
	s.user = user.Record{
		Username: "user",
		Email:    "test@cryptocades.com",
		Password: "password",
	}
	c.Assert(s.users.Create(&s.user), IsNil)

	s.user2 = user.Record{
		Username: "user2",
		Email:    "user2@cryptocades.com",
		Password: "password",
	}
	c.Assert(s.users.Create(&s.user2), IsNil)

	s.user3 = user.Record{
		Username: "user3",
		Email:    "user3@cryptocades.com",
		Password: "password",
	}
	c.Assert(s.users.Create(&s.user3), IsNil)

}

func (s *DBSuite) TearDownSuite(c *C) {
	if !testing.Short() {
		query := "Truncate users CASCADE"
		_, err := s.store.sqlx.Exec(query)
		c.Assert(err, IsNil)
	}
}

func (s *DBSuite) TearDownTest(c *C) {
	_, err := s.store.redis.Do("FLUSHALL")
	c.Assert(err, IsNil)
}

func (s *DBSuite) TestKeyName(c *C) {
	key := s.store.KeyName("daily", -1)
	d := time.Now().Add(-24 * time.Hour)
	c.Check(key, Equals, fmt.Sprintf("%d-%02d-%02d", d.Year(), d.Month(), d.Day()))

	key = s.store.KeyName("daily", 0)
	d = time.Now()
	c.Check(key, Equals, fmt.Sprintf("%d-%02d-%02d", d.Year(), d.Month(), d.Day()))

}

func (s *DBSuite) TestGet(c *C) {
	var err error
	keyname := s.store.KeyName("daily", 0)

	_, err = redis.Int(s.store.redis.Do("ZADD", keyname, 8, s.user.Id))
	c.Assert(err, IsNil)
	_, err = redis.Int(s.store.redis.Do("ZADD", keyname, 9, 2))
	c.Assert(err, IsNil)
	_, err = redis.Int(s.store.redis.Do("ZADD", keyname, 3, 3))
	c.Assert(err, IsNil)

	record, err := s.store.Get("daily", 0, 1)
	c.Assert(err, IsNil)
	c.Check(record.UserId, Equals, s.user.Id)
	c.Check(record.Score, Equals, 8)
	c.Check(record.Rank, Equals, 2)
	c.Check(record.Username, Equals, s.user.Username)

	// try a user that doesn't have any recent activity
	record, err = s.store.Get("daily", 0, 9)
	c.Assert(err, IsNil)
	c.Check(record.Score, Equals, 0)
	c.Check(record.Rank, Equals, 0)

}

func (s *DBSuite) TestTopPerformers(c *C) {
	var err error
	keyname := s.store.KeyName("daily", 0)

	_, err = redis.Int(s.store.redis.Do("ZADD", keyname, 8, s.user.Id))
	c.Assert(err, IsNil)
	_, err = redis.Int(s.store.redis.Do("ZADD", keyname, 9, s.user2.Id))
	c.Assert(err, IsNil)
	_, err = redis.Int(s.store.redis.Do("ZADD", keyname, 3, s.user3.Id))
	c.Assert(err, IsNil)

	records, err := s.store.GetTopPerformers("daily", 0, 20)
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 3)
	c.Check(records[0].UserId, Equals, s.user2.Id)
	c.Check(records[0].Score, Equals, 9)
	c.Check(records[0].Rank, Equals, 1)
	c.Check(records[0].Username, Equals, s.user2.Username)

	c.Check(records[1].UserId, Equals, s.user.Id)
	c.Check(records[1].Score, Equals, 8)
	c.Check(records[1].Rank, Equals, 2)
	c.Check(records[1].Username, Equals, s.user.Username)

	c.Check(records[2].UserId, Equals, s.user3.Id)
	c.Check(records[2].Score, Equals, 3)
	c.Check(records[2].Rank, Equals, 3)
	c.Check(records[2].Username, Equals, s.user3.Username)

}
