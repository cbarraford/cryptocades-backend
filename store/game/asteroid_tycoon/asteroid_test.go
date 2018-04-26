package asteroid_tycoon

import (
	"fmt"
	"time"

	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/cbarraford/cryptocades-backend/test"
	. "gopkg.in/check.v1"
)

type AsteroidSuite struct {
	TycoonSuite
	account Account
	ship    Ship
}

var _ = Suite(&AsteroidSuite{})

func (s *AsteroidSuite) SetUpSuite(c *C) {
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
		Health:    100,
		DrillBit:  1000,
	}
	c.Assert(s.store.CreateShip(&s.ship), IsNil)
}

func (s *AsteroidSuite) TearDownTest(c *C) {
	_, err := s.store.sqlx.Exec(fmt.Sprintf("DELETE FROM %s", asteroidsTable))
	c.Assert(err, IsNil)
}

func (s *AsteroidSuite) TestCreateAsteroid(c *C) {
	c.Assert(s.store.CreateAsteroid(&Asteroid{}), IsNil)
}

func (s *AsteroidSuite) TestMined(c *C) {
	sessionId := "abcde"
	ast := Asteroid{Total: 500, Remaining: 500, UpdatedTime: time.Now().Add(-5 * time.Second)}

	c.Assert(s.store.CreateAsteroid(&ast), IsNil)
	ship, err := s.store.GetShip(s.ship.Id)
	c.Assert(err, IsNil)
	ship.SessionId = sessionId
	c.Assert(s.store.UpdateShip(&ship), IsNil)
	c.Assert(s.store.AssignAsteroid(ast.Id, ship), IsNil)

	tx, err := s.store.sqlx.Beginx()
	c.Assert(err, IsNil)

	c.Assert(s.store.Mined(sessionId, 1, s.user.Id, tx), IsNil)
	c.Assert(tx.Commit(), IsNil)

	ast, err = s.store.OwnedAsteroid(s.ship.Id)
	c.Assert(err, IsNil)
	c.Check(ast.Remaining, Equals, 500-(1*ResourceToShareRatio))

	ship, err = s.store.GetShip(s.ship.Id)
	c.Assert(err, IsNil)
	//c.Check(ship.Health, Equals, 50)
	//c.Check(ship.DrillBit, Equals, 900)

	tx, err = s.store.sqlx.Beginx()
	defer tx.Commit()
	c.Assert(err, IsNil)
	c.Assert(s.store.Mined(sessionId, 10000, s.user.Id, tx), IsNil)
	//c.Assert(s.store.Mined(sessionId, 1, s.user.Id, tx), NotNil)
}

func (s *AsteroidSuite) TestAssign(c *C) {
	ast := Asteroid{Total: 500}
	c.Assert(s.store.CreateAsteroid(&ast), IsNil)

	ship, err := s.store.GetShip(s.ship.Id)
	c.Assert(err, IsNil)
	c.Assert(s.store.AssignAsteroid(ast.Id, ship), IsNil)
	ship.Cargo = 0
	c.Assert(s.store.AssignAsteroid(ast.Id, ship), ErrorMatches, "This asteroid is too large for your cargo hold.")
}

func (s *AsteroidSuite) TestAvailableAsteroids(c *C) {
	var err error
	ast := Asteroid{
		Remaining: 20,
	}
	c.Assert(s.store.CreateAsteroid(&ast), IsNil)

	ast2 := Asteroid{
		Total:     500,
		Remaining: 0,
	}
	c.Assert(s.store.CreateAsteroid(&ast2), IsNil)
	ship, err := s.store.GetShip(s.ship.Id)
	c.Assert(err, IsNil)
	c.Assert(s.store.AssignAsteroid(ast2.Id, ship), IsNil)

	asts, err := s.store.AvailableAsteroids()
	c.Assert(err, IsNil)
	c.Assert(asts, HasLen, 1)
	c.Assert(asts[0].Id, Equals, ast.Id)

	ast, err = s.store.OwnedAsteroid(s.ship.Id)
	c.Assert(err, IsNil)
	c.Assert(ast.Id, Equals, ast2.Id)
	c.Assert(ast.ShipSpeed, Equals, 100)
}

/*
func (s *AsteroidSuite) TestDelete(c *C) {
	var err error
	ast := Asteroid{
		Remaining: 20,
	}
	c.Assert(s.store.CreateAsteroid(&ast), IsNil)

	ast2 := Asteroid{
		Remaining: 0,
	}
	c.Assert(s.store.CreateAsteroid(&ast2), IsNil)

	c.Assert(s.store.DestroyAsteroids(), IsNil)
	asts, err := s.store.AvailableAsteroids()
	c.Assert(err, IsNil)
	c.Assert(asts, HasLen, 1)
	c.Check(asts[0].Remaining, Equals, 20)
}*/
