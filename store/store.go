package store

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"

	"github.com/cbarraford/cryptocades-backend/store/confirmation"
	"github.com/cbarraford/cryptocades-backend/store/entry"
	"github.com/cbarraford/cryptocades-backend/store/game"
	"github.com/cbarraford/cryptocades-backend/store/jackpot"
	"github.com/cbarraford/cryptocades-backend/store/session"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

type Store struct {
	Users         user.Store
	Sessions      session.Store
	Jackpots      jackpot.Store
	Entries       entry.Store
	Confirmations confirmation.Store
	Games         game.Store
}

// Get a database connection
func GetDB(url string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", url)
}

func GetRedis(url string) (redis.Conn, error) {
	return redis.DialURL(url)
}

// Get a Store object from DB connection
func GetStore(db *sqlx.DB, red redis.Conn) Store {
	return Store{
		Users:         user.NewStore(db, red),
		Sessions:      session.NewStore(db),
		Jackpots:      jackpot.NewStore(db),
		Entries:       entry.NewStore(db),
		Confirmations: confirmation.NewStore(db),
		Games:         game.NewStore(),
	}
}
