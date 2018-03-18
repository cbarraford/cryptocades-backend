package admin

import "github.com/jmoiron/sqlx"

type Store interface {
	TotalRegisterUsers() (int, error)
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
