package store

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"

	admin "github.com/cbarraford/cryptocades-backend/admin"
	"github.com/cbarraford/cryptocades-backend/store/boost"
	"github.com/cbarraford/cryptocades-backend/store/confirmation"
	"github.com/cbarraford/cryptocades-backend/store/entry"
	"github.com/cbarraford/cryptocades-backend/store/game"
	"github.com/cbarraford/cryptocades-backend/store/income"
	"github.com/cbarraford/cryptocades-backend/store/jackpot"
	"github.com/cbarraford/cryptocades-backend/store/matchup"
	"github.com/cbarraford/cryptocades-backend/store/session"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

type Store struct {
	Users         user.Store
	Sessions      session.Store
	Jackpots      jackpot.Store
	Incomes       income.Store
	Entries       entry.Store
	Confirmations confirmation.Store
	Games         game.Store
	Admins        admin.Store
	Boosts        boost.Store
	Matchups      matchup.Store
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
		Users:         user.NewStore(db),
		Sessions:      session.NewStore(db),
		Jackpots:      jackpot.NewStore(db),
		Incomes:       income.NewStore(db, red),
		Entries:       entry.NewStore(db),
		Confirmations: confirmation.NewStore(db),
		Games:         game.NewStore(),
		Admins:        admin.NewStore(db),
		Boosts:        boost.NewStore(db),
		Matchups:      matchup.NewStore(db, red),
	}
}
