package store

import (
	"github.com/jmoiron/sqlx"
)

type Store struct {
}

// Get a database connection
func GetDB(url string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", url)
}

// Get a Store object from DB connection
func GetStore(db *sqlx.DB) Store {
	return Store{}
}
