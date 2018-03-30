package asteroid_tycoon

import (
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

func (s *AsteroidSuite) SetUpTest(c *C) {
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

func (s *AsteroidSuite) TestCreateAsteroid(c *C) {
	c.Assert(s.store.CreateAsteroid(&Asteroid{}), IsNil)
}

func (s *AsteroidSuite) TestAssign(c *C) {
	ast := Asteroid{}
	c.Assert(s.store.CreateAsteroid(&ast), IsNil)

	c.Assert(s.store.AssignAsteroid(ast.Id, s.ship.Id), IsNil)
}

func (s *AsteroidSuite) TestAvailableAsteroids(c *C) {
	var err error
	ast := Asteroid{
		Remaining: 20,
	}
	c.Assert(s.store.CreateAsteroid(&ast), IsNil)

	ast2 := Asteroid{
		Remaining: 0,
	}
	c.Assert(s.store.CreateAsteroid(&ast2), IsNil)
	c.Assert(s.store.AssignAsteroid(ast2.Id, 1), IsNil)

	asts, err := s.store.AvailableAsteroids()
	c.Assert(err, IsNil)
	c.Assert(asts, HasLen, 1)
	c.Assert(asts[0].Id, Equals, ast.Id)

	asts, err = s.store.OwnedAsteroids(1)
	c.Assert(err, IsNil)
	c.Assert(asts, HasLen, 1)
	c.Assert(asts[0].Id, Equals, ast2.Id)
}

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
}
