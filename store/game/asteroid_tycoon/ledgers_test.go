package asteroid_tycoon

import (
	"time"

	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/cbarraford/cryptocades-backend/test"
	. "gopkg.in/check.v1"
)

type LedgerSuite struct {
	TycoonSuite
	account Account
	ship    Ship
}

var _ = Suite(&LedgerSuite{})

func (s *LedgerSuite) SetUpSuite(c *C) {
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

func (s *LedgerSuite) TestCompletedAsteroid(c *C) {
	var err error
	var balance int

	ast := Asteroid{
		Total:       500,
		ShipId:      s.ship.Id,
		UpdatedTime: time.Now().Add(-100000 * time.Hour),
	}

	c.Assert(s.store.CreateAsteroid(&ast), IsNil)
	ship, err := s.store.GetShip(s.ship.Id)
	c.Assert(err, IsNil)
	c.Check(ship.TotalAsteroids, Equals, 0)
	c.Check(ship.TotalResources, Equals, 0)
	c.Assert(s.store.AssignAsteroid(ast.Id, ship), IsNil)

	balance, err = s.store.ResourceBalance(s.account.Id)
	c.Assert(err, IsNil)
	c.Check(balance, Equals, 0)

	ast.Remaining = 0
	c.Assert(err, IsNil)
	c.Assert(s.store.CompletedAsteroid(ast), IsNil)

	balance, err = s.store.ResourceBalance(s.account.Id)
	c.Assert(err, IsNil)
	c.Check(balance, Equals, ast.Total)

	ship, err = s.store.GetShip(s.ship.Id)
	c.Assert(err, IsNil)
	c.Check(ship.TotalAsteroids, Equals, 1)
	c.Check(ship.TotalResources, Equals, ast.Total)
}
