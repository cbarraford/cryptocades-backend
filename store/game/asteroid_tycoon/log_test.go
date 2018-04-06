package asteroid_tycoon

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/cbarraford/cryptocades-backend/test"
)

type LogSuite struct {
	store   store
	users   user.Store
	user    user.Record
	account Account
	ship    Ship
}

var _ = Suite(&LogSuite{})

func (s *LogSuite) SetUpTest(c *C) {
	db := test.EphemeralPostgresStore(c)
	s.store = store{sqlx: db}
	s.users = user.NewStore(db)
	s.user = user.Record{
		Username: "PlayerOne",
		Email:    "playerone@cryptocades.com",
		Password: "password",
	}
	c.Assert(s.users.Create(&s.user), IsNil)

	s.account = Account{
		UserId: s.user.Id,
	}
	c.Assert(s.store.CreateAccount(&s.account), IsNil)

	s.ship = Ship{
		AccountId: s.account.Id,
	}
	c.Assert(s.store.CreateShip(&s.ship), IsNil)
}

func (s *LogSuite) TearDownSuite(c *C) {
	if !testing.Short() {
		_, err := s.store.sqlx.Exec("Truncate users CASCADE")
		c.Assert(err, IsNil)
	}
}

func (s *LogSuite) TestCreateLogRequirements(c *C) {
	line := Log{
		Log: "testing",
	}
	c.Assert(s.store.CreateLog(&line), NotNil)

	line = Log{
		ShipId: s.ship.Id,
	}
	c.Assert(s.store.CreateLog(&line), NotNil)
}

func (s *LogSuite) TestCreateLog(c *C) {
	line := Log{
		ShipId: s.ship.Id,
		Level:  2,
		Log:    "testing",
	}
	c.Assert(s.store.CreateLog(&line), IsNil)
}

func (s *LogSuite) TestGetShipLog(c *C) {
	line := Log{
		ShipId: s.ship.Id,
		Level:  2,
		Log:    "testing",
	}
	c.Assert(s.store.CreateLog(&line), IsNil)

	line = Log{
		ShipId: s.ship.Id + 1,
		Level:  2,
		Log:    "testing",
	}
	c.Assert(s.store.CreateLog(&line), IsNil)

	lines, err := s.store.GetShipLogs(s.ship.Id)
	c.Assert(err, IsNil)
	c.Assert(lines, HasLen, 1)
	c.Check(lines[0].ShipId, Equals, s.ship.Id)
}

func (s *LogSuite) TestDeleteLogs(c *C) {
	line := Log{
		ShipId: s.ship.Id,
		Level:  2,
		Log:    "testing",
	}
	c.Assert(s.store.CreateLog(&line), IsNil)

	line = Log{
		ShipId:      s.ship.Id + 1,
		Level:       2,
		Log:         "testing",
		CreatedTime: time.Now().Add(-200 * time.Hour),
	}
	c.Assert(s.store.CreateLog(&line), IsNil)

	c.Assert(s.store.DeleteOldLogs(), IsNil)

	lines, err := s.store.GetShipLogs(s.ship.Id)
	c.Assert(err, IsNil)
	c.Assert(lines, HasLen, 1)
	c.Check(lines[0].ShipId, Equals, s.ship.Id)
}
