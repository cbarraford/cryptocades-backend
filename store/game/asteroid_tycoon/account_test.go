package asteroid_tycoon

import (
	. "gopkg.in/check.v1"
)

type AccountSuite struct {
	TycoonSuite
}

var _ = Suite(&AccountSuite{})

func (s *AccountSuite) TestCreateRequirements(c *C) {
	acct := Account{}
	c.Assert(s.store.CreateAccount(&acct), NotNil)

	// ensure a user cannot have two accounts
	acct.UserId = s.user.Id
	acct.Resources = 1000
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
	acct.Resources = 300
	c.Assert(s.store.UpdateAccount(&acct), IsNil)
	acct, err = s.store.GetAccountByUserId(s.user.Id)
	c.Assert(err, IsNil)
	c.Check(acct.Credits, Equals, 100)
	c.Check(acct.Resources, Equals, 300)
	c.Check(originalUpdateTime.UnixNano(), Not(Equals), acct.UpdatedTime.UnixNano())

	// ensure we can't go negative
	acct.Credits = -100
	c.Assert(s.store.UpdateAccount(&acct), NotNil)
	acct.Credits = 100
	acct.Resources = -400
	c.Assert(s.store.UpdateAccount(&acct), NotNil)

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
