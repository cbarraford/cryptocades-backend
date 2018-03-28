package asteroid_tycoon

import (
	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/cbarraford/cryptocades-backend/test"
)

type ShipSuite struct {
	store   store
	users   user.Store
	user    user.Record
	account Account
}

var _ = Suite(&ShipSuite{})

func (s *ShipSuite) SetUpTest(c *C) {
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
}

func (s *ShipSuite) TearDownSuite(c *C) {
	_, err := s.store.sqlx.Exec("Truncate users CASCADE")
	c.Assert(err, IsNil)
}

func (s *ShipSuite) TearDownTest(c *C) {
}

func (s *ShipSuite) TestCreateRequirements(c *C) {
	ship := Ship{}
	c.Assert(s.store.CreateShip(&ship), NotNil)
}

func (s *ShipSuite) TestCreateShip(c *C) {
	ship := Ship{AccountId: s.account.Id}
	c.Assert(s.store.CreateShip(&ship), IsNil)
	c.Check(ship.Name, Equals, "Eros 433")
}

func (s *ShipSuite) TestGetShipsByUser(c *C) {
	ship := Ship{AccountId: s.account.Id}
	c.Assert(s.store.CreateShip(&ship), IsNil)

	ships, err := s.store.GetShipsByAccountId(s.account.Id)
	c.Assert(err, IsNil)
	c.Assert(ships, HasLen, 1)
	c.Check(ships[0].AccountId, Equals, int64(s.account.Id))
}

func (s *ShipSuite) TestUpdate(c *C) {
	var err error
	ship := Ship{AccountId: s.account.Id}
	c.Assert(s.store.CreateShip(&ship), IsNil)

	originalUpdateTime := ship.UpdatedTime

	ship.Name = "Forerunner"
	ship.State = 2
	ship.TotalAsteroids = 4
	ship.TotalResources = 50987
	ship.Health = 45
	ship.DrillBit = 4478
	ship.SolarSystem = 3
	c.Assert(s.store.UpdateShip(&ship), IsNil)
	ships, err := s.store.GetShipsByAccountId(s.account.Id)
	c.Assert(err, IsNil)
	c.Assert(ships, HasLen, 1)
	ship = ships[0]
	c.Check(ship.Name, Equals, "Forerunner")
	c.Check(ship.State, Equals, 2)
	c.Check(ship.TotalAsteroids, Equals, 4)
	c.Check(ship.TotalResources, Equals, 50987)
	c.Check(ship.Health, Equals, 45)
	c.Check(ship.DrillBit, Equals, 4478)
	c.Check(ship.SolarSystem, Equals, 3)
	c.Check(originalUpdateTime.UnixNano(), Not(Equals), ship.UpdatedTime.UnixNano())
}

func (s *ShipSuite) TestAddResources(c *C) {
	var err error
	ship := Ship{AccountId: s.account.Id, TotalAsteroids: 3, TotalResources: 445}
	c.Assert(s.store.CreateShip(&ship), IsNil)

	c.Assert(s.store.AddResources(2, 44), IsNil)
	ships, err := s.store.GetShipsByAccountId(s.account.Id)
	c.Assert(err, IsNil)
	c.Assert(ships, HasLen, 1)
	ship = ships[0]
	c.Check(ship.TotalAsteroids, Equals, 5)
	c.Check(ship.TotalResources, Equals, 489)
}

func (s *ShipSuite) TestAddDamage(c *C) {
	var err error
	ship := Ship{AccountId: s.account.Id, Health: 100, DrillBit: 500}
	c.Assert(s.store.CreateShip(&ship), IsNil)

	c.Assert(s.store.AddDamage(40, 60), IsNil)
	ships, err := s.store.GetShipsByAccountId(s.account.Id)
	c.Assert(err, IsNil)
	c.Assert(ships, HasLen, 1)
	ship = ships[0]
	c.Check(ship.Health, Equals, 60)
	c.Check(ship.DrillBit, Equals, 440)
}

func (s *ShipSuite) TestDelete(c *C) {
	var err error
	ship := Ship{AccountId: s.account.Id}
	c.Assert(s.store.CreateShip(&ship), IsNil)

	ships, err := s.store.GetShipsByAccountId(s.account.Id)
	c.Assert(err, IsNil)
	c.Assert(ships, HasLen, 1)
	c.Check(ships[0].AccountId, Equals, int64(s.account.Id))

	c.Assert(s.store.DeleteShip(ship.Id), IsNil)
	ships, err = s.store.GetShipsByAccountId(s.account.Id)
	c.Assert(ships, HasLen, 0)
}
