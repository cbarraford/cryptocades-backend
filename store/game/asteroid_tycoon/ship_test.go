package asteroid_tycoon

import (
	"fmt"
	"testing"
	"time"

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

func (s *ShipSuite) SetUpSuite(c *C) {
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
	if !testing.Short() {
		_, err := s.store.sqlx.Exec("Truncate users CASCADE")
		c.Assert(err, IsNil)
	}
}

func (s *ShipSuite) TearDownTest(c *C) {
	_, err := s.store.sqlx.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", shipsTable))
	c.Assert(err, IsNil)
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

func (s *ShipSuite) TestGetShipUserId(c *C) {
	ship := Ship{AccountId: s.account.Id}
	c.Assert(s.store.CreateShip(&ship), IsNil)

	userId, err := s.store.GetShipUserId(ship.Id)
	c.Assert(err, IsNil)
	c.Check(userId, Equals, s.user.Id)
}

func (s *ShipSuite) TestGetShip(c *C) {
	ship := Ship{AccountId: s.account.Id}
	c.Assert(s.store.CreateShip(&ship), IsNil)

	ship, err := s.store.GetShip(ship.Id)
	c.Assert(err, IsNil)
	c.Check(ship.Speed, Equals, 100)
	c.Check(ship.Cargo, Equals, 500)
	c.Check(ship.Repair, Equals, 1)
	c.Check(ship.Hull, Equals, 200)
}

func (s *ShipSuite) TestUpdate(c *C) {
	var err error
	ship := Ship{AccountId: s.account.Id}
	c.Assert(s.store.CreateShip(&ship), IsNil)

	originalUpdateTime := ship.UpdatedTime

	ship.Name = "Forerunner"
	ship.TotalAsteroids = 4
	ship.TotalResources = 50987
	ship.Health = 45
	ship.SolarSystem = 3
	ship.SessionId = "boo boo"
	c.Assert(s.store.UpdateShip(&ship), IsNil)
	ships, err := s.store.GetShipsByAccountId(s.account.Id)
	c.Assert(err, IsNil)
	c.Assert(ships, HasLen, 1)
	ship = ships[0]
	c.Check(ship.Name, Equals, "Forerunner")
	c.Check(ship.TotalAsteroids, Equals, 4)
	c.Check(ship.TotalResources, Equals, 50987)
	c.Check(ship.Health, Equals, 45)
	c.Check(ship.SolarSystem, Equals, 3)
	c.Check(ship.SessionId, Equals, "boo boo")
	c.Check(originalUpdateTime.UnixNano(), Not(Equals), ship.UpdatedTime.UnixNano())
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

func (s *ShipSuite) TestGetStatus(c *C) {
	status := s.store.GetStatus(Ship{}, Asteroid{})
	c.Assert(status.Status, Equals, "Docked")

	ast := Asteroid{Id: 1, Distance: 1000, ShipSpeed: 100, Total: 5000, Remaining: 5000, UpdatedTime: time.Now()}
	status = s.store.GetStatus(Ship{Id: 1, Health: 100}, ast)
	c.Assert(status.Status, Equals, "Approaching Asteroid", Commentf("Status: %+v", status))
	c.Assert(status.RemainingTime, Equals, 10)

	ast = Asteroid{Id: 1, Distance: 1000, ShipSpeed: 100, Total: 5000, Remaining: 5000, UpdatedTime: time.Now().Add(-100 * time.Second)}
	status = s.store.GetStatus(Ship{Id: 1, Health: 100}, ast)
	c.Assert(status.Status, Equals, "Mining", Commentf("Status: %+v", status))
	c.Assert(status.RemainingTime, Equals, -90)

	ast = Asteroid{Id: 1, Distance: 2000, ShipSpeed: 100, Total: 5000, Remaining: 400, UpdatedTime: time.Now()}
	status = s.store.GetStatus(Ship{Id: 1, Health: 100}, ast)
	c.Assert(status.Status, Equals, "Mining")
	c.Assert(status.RemainingTime, Equals, 20)

	ast = Asteroid{Id: 1, Distance: 2000, ShipSpeed: 100, Total: 5000, Remaining: 0, UpdatedTime: time.Now()}
	status = s.store.GetStatus(Ship{Id: 1, Health: 100}, ast)
	c.Assert(status.Status, Equals, "Approaching Space Station")
	c.Assert(status.RemainingTime, Equals, 20)

	ast = Asteroid{Id: 1, Distance: 2000, ShipSpeed: 100, Total: 5000, Remaining: 0, UpdatedTime: time.Now().Add(-100 * time.Second)}
	status = s.store.GetStatus(Ship{Id: 1, Health: 100}, ast)
	c.Assert(status.Status, Equals, "Docked")
	c.Assert(status.RemainingTime, Equals, -80)
}
