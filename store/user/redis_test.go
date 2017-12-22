package user

import (
	"fmt"

	. "gopkg.in/check.v1"

	"github.com/CBarraford/lotto/test"
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
	record := Record{
		Username: "bob",
		Email:    "bob@lotto.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	for i, member := range []string{"red", "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq", "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq"} {
		s.store.redis.Send("ZADD", "hashes", i, member)
	}
	if _, err := s.store.redis.Do(""); err != nil {
		c.Assert(err, IsNil)
	}

	c.Assert(s.store.UpdateScores(), IsNil)
	c.Assert(s.store.UpdateScores(), IsNil)

	record, err = s.store.Get(record.Id)
	c.Assert(err, IsNil)
	c.Check(record.MinedHashes, Equals, 2)
}