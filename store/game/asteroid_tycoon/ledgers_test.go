package asteroid_tycoon

import (
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

func (s *LedgerSuite) SetUpTest(c *C) {
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
		ShipId: s.ship.Id,
	}

	c.Assert(s.store.CreateAsteroid(&ast), IsNil)
	c.Assert(s.store.AssignAsteroid(ast.Id, s.ship.Id), IsNil)

	balance, err = s.store.ResourceBalance(s.account.Id)
	c.Assert(err, IsNil)
	c.Check(balance, Equals, 0)

	c.Assert(s.store.CompletedAsteroid(ast), IsNil)

	balance, err = s.store.ResourceBalance(s.account.Id)
	c.Assert(err, IsNil)
	c.Check(balance, Equals, ast.Total)

}
