package context

import (
	"fmt"
	"go/build"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"

	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store"
	"github.com/cbarraford/cryptocades-backend/util"
)

func TestPackage(t *testing.T) { TestingT(t) }

type ContextSuite struct {
	store store.Store
}

var _ = Suite(&ContextSuite{})

func (s *ContextSuite) TestMigrate(c *C) {
	if testing.Short() {
		c.Skip("Short mode: no integration tests")
	}

	ci := os.Getenv("CI")
	var dbURL string
	if ci == "1" {
		dbURL = "postgres://ubuntu@localhost:5432/test?sslmode=disable"
	} else {
		dbURL = "postgres://postgres:password@postgres:5432/db?sslmode=disable"
	}
	dbx, err := sqlx.Connect("postgres", dbURL)
	c.Assert(err, IsNil)

	// create database and select
	dbname := util.RandSeq(16, util.LowerLetters) // databases must be lower case
	_ = sqlx.MustExec(dbx, fmt.Sprintf("CREATE DATABASE %s;", dbname))

	dbx.Close()

	if ci == "1" {
		dbURL = fmt.Sprintf("postgres://ubuntu@localhost:5432/%s?sslmode=disable", dbname)
	} else {
		dbURL = fmt.Sprintf("postgres://postgres:password@postgres:5432/%s?sslmode=disable", dbname)
	}

	migrateDir := fmt.Sprintf("file://%s/src/github.com/cbarraford/cryptocades-backend/migrations", build.Default.GOPATH)
	c.Assert(MigrateDB(dbURL, migrateDir), IsNil)

	// do a second time so that no change does not create an error
	c.Assert(MigrateDB(dbURL, migrateDir), IsNil)
}
