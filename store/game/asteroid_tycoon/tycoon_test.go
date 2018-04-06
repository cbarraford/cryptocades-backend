package asteroid_tycoon

import (
	"testing"

	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/cbarraford/cryptocades-backend/test"
)

func TestPackage(t *testing.T) { TestingT(t) }

type TycoonSuite struct {
	store store
	users user.Store
	user  user.Record
}

func (s *TycoonSuite) SetUpTest(c *C) {
	db := test.EphemeralPostgresStore(c)
	s.store = store{sqlx: db}
	s.users = user.NewStore(db)
	s.user = user.Record{
		Username: "PlayerOne",
		Email:    "playerone@cryptocades.com",
		Password: "password",
	}
	c.Assert(s.users.Create(&s.user), IsNil)
}

func (s *TycoonSuite) TearDownSuite(c *C) {
	if !testing.Short() {
		_, err := s.store.sqlx.Exec("Truncate users CASCADE")
		c.Assert(err, IsNil)
	}
}
