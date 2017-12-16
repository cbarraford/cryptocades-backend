package test

import (
	"fmt"
	"os"
	"testing"

	"go/build"

	check "gopkg.in/check.v1"

	"github.com/jmoiron/sqlx"

	"github.com/CBarraford/lotto/store/context"
	"github.com/CBarraford/lotto/util"
)

const NoIntegration string = "Short mode: no integration tests"

func EphemeralURLStore(c *check.C) string {
	if testing.Short() {
		c.Skip(NoIntegration)
		return ""
	}

	ci := os.Getenv("CI")
	var dbURL string
	if ci == "true" {
		dbURL = "postgres://ubuntu@localhost:5432/test"
	} else {
		dbURL = "postgres://postgres:password@postgres:5432/db?sslmode=disable"
	}
	dbx, err := sqlx.Connect("postgres", dbURL)
	c.Assert(err, check.IsNil)

	// create database and select
	dbname := util.RandSeq(16, util.LowerLetters) // databases must be lower case
	_ = sqlx.MustExec(dbx, fmt.Sprintf("CREATE DATABASE %s;", dbname))

	dbx.Close()

	if ci == "true" {
		dbURL = fmt.Sprintf("postgres://ubuntu@localhost:5432/%s?sslmode=disable", dbname)
	} else {
		dbURL = fmt.Sprintf("postgres://postgres:password@postgres:5432/%s?sslmode=disable", dbname)
	}

	migrateDir := fmt.Sprintf("file://%s/src/github.com/CBarraford/lotto/migrations", build.Default.GOPATH)
	err = context.MigrateDB(url, migrateDir)
	c.Assert(err, check.IsNil)

	return url
}

// EphemeralPostgresStore returns a connection to a randomly generated database
func EphemeralPostgresStore(c *check.C) *sqlx.DB {
	url := EphemeralURLStore(c)

	db, err := sqlx.Connect("postgres", url)
	c.Assert(err, check.IsNil)
	return db
}
