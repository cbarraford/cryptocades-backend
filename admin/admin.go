package admin

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	TotalRegisterUsers() (int, error)
	TotalActiveUsers(minutes int) (int, error)
	AwardPlays(email string, amount int, description string) error
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
		fmt.Sprintf("SELECT COUNT(id) FROM incomes WHERE game_id > 0 AND updated_time >= (now() - '%d minute'::INTERVAL) GROUP BY user_id;", minutes),
	)
	err = db.sqlx.Get(&i, query)
	return
}

func (db *store) AwardPlays(email string, amount int, title string) error {
	var userId int64
	query := db.sqlx.Rebind("SELECT id FROM users WHERE email = ?")
	err := db.sqlx.Get(&userId, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Email not found")
		} else {
			return err
		}
	}

	query = db.sqlx.Rebind("INSERT INTO incomes (user_id, game_id, session_id, amount, partial_amount) VALUES (?,0,?,?,0)")
	_, err = db.sqlx.Exec(query, userId, title, amount)
	return err
}
