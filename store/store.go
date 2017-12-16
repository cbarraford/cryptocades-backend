package store

import (
	"github.com/jmoiron/sqlx"

	"github.com/CBarraford/lotto/store/user"
)

type Store struct {
	Users user.Store
}

// Get a database connection
func GetDB(url string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", url)
}

// Get a Store object from DB connection
func GetStore(db *sqlx.DB) Store {
	return Store{
		Users: user.NewStore(db),
	}
}
