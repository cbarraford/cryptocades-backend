package admin

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	TotalRegisterUsers() (int, error)
	TotalActiveUsers(minutes int) (int, error)
}

type store struct {
	Store
	sqlx *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return &store{sqlx: db}
}

func (db *store) TotalRegisterUsers() (i int, err error) {
	query := "SELECT COUNT(id) FROM users"
	err = db.sqlx.Get(&i, query)
	return
}

func (db *store) TotalActiveUsers(minutes int) (i int, err error) {
	query := db.sqlx.Rebind(
		fmt.Sprintf("SELECT COUNT(id) FROM incomes WHERE game_id > 0 AND updated_time >= (now() - '%d minute'::INTERVAL);", minutes),
	)
	err = db.sqlx.Get(&i, query)
	return
}
