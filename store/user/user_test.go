package user

import (
	"fmt"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/test"
)

func TestPackage(t *testing.T) { TestingT(t) }

type DBSuite struct {
	store store
}

var _ = Suite(&DBSuite{})

func (s *DBSuite) SetUpSuite(c *C) {
	db := test.EphemeralPostgresStore(c)
	s.store = store{sqlx: db}
}

func (s *DBSuite) TearDownTest(c *C) {
	query := fmt.Sprintf("Truncate %s CASCADE", table)
	_, err := s.store.sqlx.Exec(query)
	c.Assert(err, IsNil)
}

func (s *DBSuite) TestCreateRequirements(c *C) {
	record := Record{
		Username: "bob",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), NotNil)

	record = Record{
		Email:    "bob@cryptocades.com",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), NotNil)

	record = Record{
		Email:    "bob@cryptocades.com",
		Username: "bob",
	}
	c.Assert(s.store.Create(&record), NotNil)
}

func (s *DBSuite) TestCreate(c *C) {
	record := Record{
		Username: "bob",
		Email:    "bob@cryptocades.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)
	c.Check(record.Username, Equals, "bob")
	c.Check(record.Email, Equals, "bob@cryptocades.com")
	c.Check(record.BTCAddr, Equals, "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq")
	c.Check(CheckPasswordHash("password", record.Password), Equals, true)
}

func (s *DBSuite) TestGet(c *C) {
	record := Record{
		Username: "bob",
		Email:    "bob@cryptocades.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	r, err := s.store.Get(record.Id)
	c.Assert(err, IsNil)
	c.Check(r.Username, Equals, "bob")
	c.Check(r.Email, Equals, "bob@cryptocades.com")
	c.Check(r.BTCAddr, Equals, "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq")
	c.Check(CheckPasswordHash("password", r.Password), Equals, true)
}

func (s *DBSuite) TestGetByUsername(c *C) {
	record := Record{
		Username: "bob",
		Email:    "bob@cryptocades.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	r, err := s.store.GetByUsername(record.Username)
	c.Assert(err, IsNil)
	c.Check(r.Username, Equals, "bob")
	c.Check(r.Email, Equals, "bob@cryptocades.com")
	c.Check(r.BTCAddr, Equals, "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq")
	c.Check(CheckPasswordHash("password", r.Password), Equals, true)
}

func (s *DBSuite) TestGetByBTCAddress(c *C) {
	record := Record{
		Username: "bob",
		Email:    "bob@cryptocades.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	r, err := s.store.GetByBTCAddress(record.BTCAddr)
	c.Assert(err, IsNil)
	c.Check(r.Username, Equals, "bob")
	c.Check(r.Email, Equals, "bob@cryptocades.com")
	c.Check(r.BTCAddr, Equals, "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq")
	c.Check(CheckPasswordHash("password", r.Password), Equals, true)
}

func (s *DBSuite) TestGetByEmail(c *C) {
	record := Record{
		Username: "bob",
		Email:    "bob@cryptocades.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	r, err := s.store.GetByEmail(record.Email)
	c.Assert(err, IsNil)
	c.Check(r.Username, Equals, "bob")
	c.Check(r.Email, Equals, "bob@cryptocades.com")
	c.Check(r.BTCAddr, Equals, "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq")
	c.Check(CheckPasswordHash("password", r.Password), Equals, true)
}

func (s *DBSuite) TestGetByReferralCode(c *C) {
	record := Record{
		Username: "bob",
		Email:    "bob@cryptocades.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	var err error
	record, err = s.store.Get(record.Id)
	c.Assert(err, IsNil)

	r, err := s.store.GetByReferralCode(record.ReferralCode)
	c.Assert(err, IsNil)
	c.Check(r.Username, Equals, "bob")
	c.Check(r.Email, Equals, "bob@cryptocades.com")
	c.Check(r.BTCAddr, Equals, "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq")
	c.Check(r.ReferralCode, Equals, record.ReferralCode, Commentf("%s", record.ReferralCode))
	c.Check(CheckPasswordHash("password", r.Password), Equals, true)
}

func (s *DBSuite) TestUpdate(c *C) {
	record := Record{
		Username: "bob",
		Email:    "bob@cryptocades.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	r, err := s.store.GetByBTCAddress(record.BTCAddr)
	c.Assert(err, IsNil)
	r.Username = "bobby"
	r.Email = "bobby@cryptocades.com"
	r.BTCAddr = "NMiJFQvupX5kSZcUtfSoD9NtLevUgjv3ur"

	c.Assert(s.store.Update(&r), IsNil)

	r, err = s.store.GetByBTCAddress(r.BTCAddr)
	c.Assert(err, IsNil)
	c.Check(r.Username, Equals, "bobby")
	c.Check(r.Email, Equals, "bobby@cryptocades.com")
	c.Check(r.BTCAddr, Equals, "NMiJFQvupX5kSZcUtfSoD9NtLevUgjv3ur")
	c.Check(r.UpdatedTime, Not(Equals), record.UpdatedTime)
}

func (s *DBSuite) TestList(c *C) {
	record := Record{
		Username: "bob",
		Email:    "bob@cryptocades.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	records, err := s.store.List()
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 1)
	r := records[0]
	c.Check(r.Username, Equals, "bob")
	c.Check(r.Email, Equals, "bob@cryptocades.com")
	c.Check(r.BTCAddr, Equals, "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq")
	c.Check(CheckPasswordHash("password", r.Password), Equals, true)
}

func (s *DBSuite) TestAuthenticate(c *C) {
	var err error
	record := Record{
		Username: "bob",
		Email:    "bob@cryptocades.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	// failure because the user isn't confirmed yet
	_, err = s.store.Authenticate("bob", "password")
	c.Assert(err, NotNil)

	record.Email = "bob@bob.com"
	c.Assert(s.store.MarkAsConfirmed(&record), IsNil)

	// happy path
	record, err = s.store.Authenticate("bob", "password")
	c.Assert(err, IsNil)
	c.Check(record.Username, Equals, "bob")
	// ensure email gets updated when we mark at confirmed
	c.Check(record.Email, Equals, "bob@bob.com")

	// bad password
	record, err = s.store.Authenticate("bob", "bad password")
	c.Assert(err, ErrorMatches, "Incorrect username or password")
	record, err = s.store.Authenticate("bob", "")
	c.Assert(err, ErrorMatches, "Incorrect username or password")

	// bad username
	record, err = s.store.Authenticate("bad username", "bad password")
	c.Assert(err, ErrorMatches, "Incorrect username or password")
}

func (s *DBSuite) TestPasswordSet(c *C) {
	var err error
	record := Record{
		Username: "bob",
		Email:    "bob@cryptocades.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	record.Password = "another passwd"
	c.Assert(s.store.PasswordSet(&record), IsNil)

	// happy path
	record, err = s.store.Authenticate("bob", "another passwd")
	c.Assert(err, IsNil)
	c.Check(record.Username, Equals, "bob")
}

func (s *DBSuite) TestDelete(c *C) {
	var err error
	record := Record{
		Username: "bob",
		Email:    "bob@cryptocades.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	record, err = s.store.Get(record.Id)
	c.Assert(err, IsNil)

	c.Assert(s.store.Delete(record.Id), IsNil)

	record, err = s.store.Get(record.Id)
	c.Assert(err, NotNil)

}
