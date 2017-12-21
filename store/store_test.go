package store

import (
	"os"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/CBarraford/lotto/store/user"
	"github.com/CBarraford/lotto/test"
)

func TestPackage(t *testing.T) { TestingT(t) }

type StoreSuite struct{}

var _ = Suite(&StoreSuite{})

func (s *StoreSuite) TestEphemeralPostgres(c *C) {
	url := test.EphemeralURLStore(c)

	db, err := GetDB(url)
	c.Assert(err, IsNil)

	red, err := GetRedis(os.Getenv("REDIS_URL"))
	c.Assert(err, IsNil)

	cstore := GetStore(db, red)

	record := user.Record{
		Username: "bob",
		Email:    "bob@lotto.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(cstore.Users.Create(&record), IsNil)

	users, err := cstore.Users.List()
	c.Assert(err, IsNil)
	c.Assert(users, HasLen, 1)
}
