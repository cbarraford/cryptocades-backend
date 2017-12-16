package user

import (
	"fmt"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/CBarraford/lotto/test"
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
	query := fmt.Sprintf("Truncate %s", table)
	_, err := s.store.sqlx.Exec(query)
	c.Assert(err, IsNil)
}

func (s *DBSuite) TestCreateRequirements(c *C) {
	record := Record{
		Username: "bob",
		Password: "password",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
	}
	c.Assert(s.store.Create(&record), NotNil)

	record = Record{
		Email:    "bob@lotto.com",
		Password: "password",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
	}
	c.Assert(s.store.Create(&record), NotNil)

	record = Record{
		Email:    "bob@lotto.com",
		Username: "bob",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
	}
	c.Assert(s.store.Create(&record), NotNil)

	record = Record{
		Email:    "bob@lotto.com",
		Username: "bob",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), NotNil)
}

func (s *DBSuite) TestCreate(c *C) {
	record := Record{
		Username: "bob",
		Email:    "bob@lotto.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)
	c.Check(record.Username, Equals, "bob")
	c.Check(record.Email, Equals, "bob@lotto.com")
	c.Check(record.BTCAddr, Equals, "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq")
	c.Check(CheckPasswordHash("password", record.Password), Equals, true)
}

func (s *DBSuite) TestGet(c *C) {
	record := Record{
		Username: "bob",
		Email:    "bob@lotto.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	r, err := s.store.Get(record.Id)
	c.Assert(err, IsNil)
	c.Check(r.Username, Equals, "bob")
	c.Check(r.Email, Equals, "bob@lotto.com")
	c.Check(r.BTCAddr, Equals, "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq")
	c.Check(CheckPasswordHash("password", r.Password), Equals, true)
}

func (s *DBSuite) TestGetByUsername(c *C) {
	record := Record{
		Username: "bob",
		Email:    "bob@lotto.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	r, err := s.store.GetByUsername(record.Username)
	c.Assert(err, IsNil)
	c.Check(r.Username, Equals, "bob")
	c.Check(r.Email, Equals, "bob@lotto.com")
	c.Check(r.BTCAddr, Equals, "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq")
	c.Check(r.TotalHashes, Equals, 0)
	c.Check(r.BonusHashes, Equals, 0)
	c.Check(CheckPasswordHash("password", r.Password), Equals, true)
}

func (s *DBSuite) TestGetByBTCAddress(c *C) {
	record := Record{
		Username: "bob",
		Email:    "bob@lotto.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	r, err := s.store.GetByBTCAddress(record.BTCAddr)
	c.Assert(err, IsNil)
	c.Check(r.Username, Equals, "bob")
	c.Check(r.Email, Equals, "bob@lotto.com")
	c.Check(r.BTCAddr, Equals, "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq")
	c.Check(CheckPasswordHash("password", r.Password), Equals, true)
}

func (s *DBSuite) TestUpdate(c *C) {
	record := Record{
		Username: "bob",
		Email:    "bob@lotto.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	r, err := s.store.GetByBTCAddress(record.BTCAddr)
	c.Assert(err, IsNil)
	r.Username = "bobby"
	r.Email = "bobby@lotto.com"
	r.TotalHashes = 101
	r.BonusHashes = 10
	r.BTCAddr = "NMiJFQvupX5kSZcUtfSoD9NtLevUgjv3ur"

	c.Assert(s.store.Update(&r), IsNil)

	r, err = s.store.GetByBTCAddress(r.BTCAddr)
	c.Assert(err, IsNil)
	c.Check(r.Username, Equals, "bobby")
	c.Check(r.Email, Equals, "bobby@lotto.com")
	c.Check(r.TotalHashes, Equals, 101)
	c.Check(r.BonusHashes, Equals, 10)
	c.Check(r.BTCAddr, Equals, "NMiJFQvupX5kSZcUtfSoD9NtLevUgjv3ur")
	c.Check(r.UpdatedTime, Not(Equals), record.UpdatedTime)
}

func (s *DBSuite) TestList(c *C) {
	record := Record{
		Username: "bob",
		Email:    "bob@lotto.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	records, err := s.store.List()
	c.Assert(err, IsNil)
	c.Assert(records, HasLen, 1)
	r := records[0]
	c.Check(r.Username, Equals, "bob")
	c.Check(r.Email, Equals, "bob@lotto.com")
	c.Check(r.BTCAddr, Equals, "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq")
	c.Check(CheckPasswordHash("password", r.Password), Equals, true)
}

func (s *DBSuite) TestAuthenticate(c *C) {
	var err error
	record := Record{
		Username: "bob",
		Email:    "bob@lotto.com",
		BTCAddr:  "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
		Password: "password",
	}
	c.Assert(s.store.Create(&record), IsNil)

	// happy path
	record, err = s.store.Authenticate("bob", "password")
	c.Assert(err, IsNil)
	c.Check(record.Username, Equals, "bob")

	// bad password
	record, err = s.store.Authenticate("bob", "bad password")
	c.Assert(err, ErrorMatches, "Incorrect username or password")
	record, err = s.store.Authenticate("bob", "")
	c.Assert(err, ErrorMatches, "Incorrect username or password")

	// bad username
	record, err = s.store.Authenticate("bad username", "bad password")
	c.Assert(err, ErrorMatches, "Incorrect username or password")
}
