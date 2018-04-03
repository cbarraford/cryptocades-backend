package income

import (
	"fmt"

	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store/game/asteroid_tycoon"
	"github.com/cbarraford/cryptocades-backend/test"
)

type RedisSuite struct {
	store store
}

var _ = Suite(&RedisSuite{})

func (s *RedisSuite) SetUpSuite(c *C) {
	db := test.EphemeralPostgresStore(c)
	red := test.EphemeralRedisStore(c)
	s.store = store{sqlx: db, redis: red}
}

func (s *RedisSuite) TearDownTest(c *C) {
	query := fmt.Sprintf("Truncate %s", table)
	_, err := s.store.sqlx.Exec(query)
	c.Assert(err, IsNil)

	_, err = s.store.redis.Do("FLUSHALL")
	c.Assert(err, IsNil)
}

func (s *RedisSuite) TestZPop(c *C) {
	var err error
	for _, member := range []string{
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-red@1",
		"1-blue@1",
		"12@1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		"",
		"bogus",
		"1-green@12",
		"green@1",
	} {
		s.store.redis.Send("ZINCRBY", "shares", 1, member)
	}
	if _, err := s.store.redis.Do(""); err != nil {
		c.Assert(err, IsNil)
	}

	ty := asteroid_tycoon.NewStore(s.store.sqlx)
	c.Assert(s.store.UpdateScores(ty), IsNil)
	c.Assert(s.store.UpdateScores(ty), IsNil)

	total, err := s.store.UserIncome(1)
	c.Assert(err, IsNil)
	c.Check(total, Equals, 1)
}
