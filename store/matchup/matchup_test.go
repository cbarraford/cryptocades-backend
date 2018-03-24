package matchup

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

func (s *DBSuite) TestKeyName(c *C) {
	key := s.store.KeyName("daily", -1)
	d := time.now().Add(-24 * time.Hour)
	c.Check(key, Equals, fmt.Sprintf("%d-%02d-%02d", d.Year(), d.Month(), d.Day()))
}
