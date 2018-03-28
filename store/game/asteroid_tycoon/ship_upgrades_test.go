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
	upgrade := AppliedUpgrade{
		CategoryId: 2,
		AssetId:    3,
	}
	c.Assert(s.store.ApplyUpgrade(&upgrade), NotNil)

	upgrade = AppliedUpgrade{
		ShipId:  1,
		AssetId: 3,
	}
	c.Assert(s.store.ApplyUpgrade(&upgrade), NotNil)

	upgrade = AppliedUpgrade{
		ShipId:     1,
		CategoryId: 2,
	}
	c.Assert(s.store.ApplyUpgrade(&upgrade), NotNil)
}

func (s *ShipUpgradeSuite) TestApplyUpgrade(c *C) {
	ship := Ship{AccountId: s.account.Id}
	c.Assert(s.store.CreateShip(&ship), IsNil)
	upgrade := AppliedUpgrade{ShipId: ship.Id, CategoryId: 1, AssetId: 2}
	c.Assert(s.store.ApplyUpgrade(&upgrade), IsNil)
}

func (s *ShipUpgradeSuite) TestGetUpgradesByShipId(c *C) {
	ship := Ship{AccountId: s.account.Id}
	c.Assert(s.store.CreateShip(&ship), IsNil)
	ship2 := Ship{AccountId: s.account.Id}
	c.Assert(s.store.CreateShip(&ship2), IsNil)

	upgrade := AppliedUpgrade{ShipId: ship2.Id, CategoryId: 1, AssetId: 6}
	c.Assert(s.store.ApplyUpgrade(&upgrade), IsNil)
	upgrade = AppliedUpgrade{ShipId: ship.Id, CategoryId: 1, AssetId: 2}
	c.Assert(s.store.ApplyUpgrade(&upgrade), IsNil)
	upgrade = AppliedUpgrade{ShipId: ship.Id, CategoryId: 1, AssetId: 3}
	c.Assert(s.store.ApplyUpgrade(&upgrade), IsNil)
	upgrade = AppliedUpgrade{ShipId: ship.Id, CategoryId: 1, AssetId: 4}
	c.Assert(s.store.ApplyUpgrade(&upgrade), IsNil)

	upgrades, err := s.store.GetUpgradesByShipId(ship.Id)
	c.Assert(err, IsNil)
	c.Assert(upgrades, HasLen, 1)
	c.Check(upgrades[0].ShipId, Equals, ship.Id)
	c.Check(upgrades[0].CategoryId, Equals, 1)
	c.Check(upgrades[0].AssetId, Equals, 4)
}
