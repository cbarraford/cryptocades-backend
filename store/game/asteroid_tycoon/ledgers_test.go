package asteroid_tycoon

import (
	. "gopkg.in/check.v1"
)

type LedgerSuite struct {
	TycoonSuite
}

var _ = Suite(&LedgerSuite{})

func (s *LedgerSuite) TestCreateRequirements(c *C) {
	acct := Account{}
	c.Assert(s.store.CreateAccount(&acct), NotNil)

	// ensure a user cannot have two accounts
	acct.UserId = s.user.Id
	acct.Resources = 1000
	c.Assert(s.store.CreateAccount(&acct), IsNil)
	c.Assert(s.store.CreateAccount(&acct), NotNil)
}

func (s *LedgerSuite) TestCreateAccount(c *C) {
	c.Assert(s.store.CreateAccount(&Account{UserId: s.user.Id}), IsNil)
}
