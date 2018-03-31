package asteroid_tycoon

import (
	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/cbarraford/cryptocades-backend/test"
)

type ShipUpgradeSuite struct {
	TycoonSuite
	account Account
}

var _ = Suite(&ShipUpgradeSuite{})

func (s *ShipUpgradeSuite) SetUpTest(c *C) {
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

func (s *ShipUpgradeSuite) TestCreateRequirements(c *C) {
	upgrade := ShipUpgrade{
		CategoryId: 2,
		AssetId:    3,
	}
	c.Assert(s.store.ApplyUpgrade(0, upgrade), NotNil)

	upgrade = ShipUpgrade{
		AssetId: 3,
	}
	c.Assert(s.store.ApplyUpgrade(1, upgrade), NotNil)

	upgrade = ShipUpgrade{
		CategoryId: 2,
	}
	c.Assert(s.store.ApplyUpgrade(1, upgrade), NotNil)
}

func (s *ShipUpgradeSuite) TestIniShip(c *C) {
	ship := Ship{AccountId: s.account.Id}
	c.Assert(s.store.CreateShip(&ship), IsNil)

	upgrades, err := s.store.GetUpgradesByShipId(ship.Id)
	c.Assert(err, IsNil)
	c.Assert(upgrades, HasLen, 4)
	c.Assert(upgrades[0].CategoryId, Equals, 1)
	c.Assert(upgrades[0].AssetId, Equals, 1)
	c.Assert(upgrades[0].ShipId, Equals, ship.Id)
}

func (s *ShipUpgradeSuite) TestApplyUpgrade(c *C) {
	ship := Ship{AccountId: s.account.Id}
	c.Assert(s.store.CreateShip(&ship), IsNil)

	upgrade := ShipUpgrade{CategoryId: 1, AssetId: 6}
	c.Assert(s.store.ApplyUpgrade(ship.Id, upgrade), ErrorMatches, "Upgrades can only be applied in sequential order")
	upgrade = ShipUpgrade{CategoryId: 1, AssetId: 2, Cost: 100}
	c.Assert(s.store.ApplyUpgrade(ship.Id, upgrade), ErrorMatches, "Insufficient funds.")
	upgrade = ShipUpgrade{CategoryId: 1, AssetId: 2, Cost: 0}
	c.Assert(s.store.ApplyUpgrade(ship.Id, upgrade), IsNil)
}

func (s *ShipUpgradeSuite) TestGetUpgradesByShipId(c *C) {
	ship := Ship{AccountId: s.account.Id}
	c.Assert(s.store.CreateShip(&ship), IsNil)
	ship2 := Ship{AccountId: s.account.Id}
	c.Assert(s.store.CreateShip(&ship2), IsNil)

	upgrade := ShipUpgrade{CategoryId: 1, AssetId: 2}
	c.Assert(s.store.ApplyUpgrade(ship.Id, upgrade), IsNil)
	upgrade = ShipUpgrade{CategoryId: 1, AssetId: 3}
	c.Assert(s.store.ApplyUpgrade(ship.Id, upgrade), IsNil)
	upgrade = ShipUpgrade{CategoryId: 1, AssetId: 4}
	c.Assert(s.store.ApplyUpgrade(ship.Id, upgrade), IsNil)

	upgrades, err := s.store.GetUpgradesByShipId(ship.Id)
	c.Assert(err, IsNil)
	c.Assert(upgrades, HasLen, 4)
	c.Check(upgrades[0].ShipId, Equals, ship.Id)
	c.Check(upgrades[0].CategoryId, Equals, 1)
	c.Check(upgrades[0].AssetId, Equals, 4)
}
