package asteroid_tycoon

import (
	"fmt"

	. "gopkg.in/check.v1"
)

type AccountSuite struct {
	TycoonSuite
}

var _ = Suite(&AccountSuite{})

func (s *AccountSuite) TearDownTest(c *C) {
	_, err := s.store.sqlx.Exec(fmt.Sprintf("Truncate %s CASCADE", accountsTable))
	c.Assert(err, IsNil)
}

func (s *AccountSuite) TestCreateRequirements(c *C) {
	acct := Account{}
	c.Assert(s.store.CreateAccount(&acct), NotNil)

	// ensure a user cannot have two accounts
	acct.UserId = s.user.Id
	c.Assert(s.store.CreateAccount(&acct), IsNil)
	c.Assert(s.store.CreateAccount(&acct), NotNil)
}

func (s *AccountSuite) TestCreateAccount(c *C) {
	c.Assert(s.store.CreateAccount(&Account{UserId: s.user.Id}), IsNil)
}

func (s *AccountSuite) TestGetAccountByUser(c *C) {
	c.Assert(s.store.CreateAccount(&Account{UserId: s.user.Id}), IsNil)

	acct, err := s.store.GetAccountByUserId(s.user.Id)
	c.Assert(err, IsNil)
	c.Check(acct.UserId, Equals, s.user.Id)
}

func (s *AccountSuite) TestUpdate(c *C) {
	var err error
	acct := Account{UserId: s.user.Id}
	c.Assert(s.store.CreateAccount(&acct), IsNil)

	originalUpdateTime := acct.UpdatedTime

	acct.Credits = 100
	c.Assert(s.store.UpdateAccount(&acct), IsNil)
	acct, err = s.store.GetAccountByUserId(s.user.Id)
	c.Assert(err, IsNil)
	c.Check(acct.Credits, Equals, 100)
	c.Check(originalUpdateTime.UnixNano(), Not(Equals), acct.UpdatedTime.UnixNano())

	// ensure we can't go negative
	acct.Credits = -100
	c.Assert(s.store.UpdateAccount(&acct), NotNil)
}

func (s *AccountSuite) TestCredits(c *C) {
	var err error
	acct := Account{UserId: s.user.Id}
	c.Assert(s.store.CreateAccount(&acct), IsNil)

	c.Assert(s.store.AddAccountCredits(acct.Id, 100), IsNil)
	c.Assert(s.store.SubtractAccountCredits(acct.Id, 40), IsNil)

	acct, err = s.store.GetAccountByUserId(s.user.Id)
	c.Assert(err, IsNil)
	c.Check(acct.Credits, Equals, 60)
}

func (s *AccountSuite) TestDelete(c *C) {
	var err error
	acct := Account{UserId: s.user.Id}
	c.Assert(s.store.CreateAccount(&acct), IsNil)

	acct, err = s.store.GetAccountByUserId(s.user.Id)
	c.Assert(err, IsNil)
	c.Check(acct.UserId, Equals, s.user.Id)

	c.Assert(s.store.DeleteAccount(acct.Id), IsNil)
	acct, err = s.store.GetAccountByUserId(s.user.Id)
	c.Assert(err, NotNil)
}

func (s *AccountSuite) TestTradeForPlays(c *C) {
	acct := Account{UserId: s.user.Id}
	c.Assert(s.store.CreateAccount(&acct), IsNil)

	c.Assert(s.store.AddAccountCredits(acct.Id, 1000), IsNil)
	c.Assert(s.store.TradeForPlays(acct.Id, 1000), ErrorMatches, "Insufficient funds.")
	c.Assert(s.store.TradeForPlays(acct.Id, 1), IsNil)

}

func (s *AccountSuite) TestTradeForCredits(c *C) {
	acct := Account{UserId: s.user.Id}
	c.Assert(s.store.CreateAccount(&acct), IsNil)

	c.Assert(s.store.TradeForCredits(acct.Id, 1000), ErrorMatches, "Insufficient funds.")
	c.Assert(s.store.AddEntry(acct.Id, 5000, "test"), IsNil)
	c.Assert(s.store.TradeForCredits(acct.Id, 5), IsNil)
	c.Assert(s.store.TradeForCredits(acct.Id, 1), ErrorMatches, "Insufficient funds.")
}
