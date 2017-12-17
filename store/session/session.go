package session

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/CBarraford/lotto/util"
)

type Store interface {
	Create(record *Record, i int) error
	GetByToken(token string) (Record, error)
	Authenticate(token string) (int64, error)
	Delete(token string) error
}

type store struct {
	Store
	sqlx *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return &store{sqlx: db}
}

const table string = "sessions"

type Record struct {
	Id          int64     `json:"id" db:"id"`
	UserId      int64     `json:"user_id" db:"user_id"`
	Token       string    `json:"token" db:"token"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
	ExpireTime  time.Time `json:"expire_time" db:"expire_time"`
}

func (db *store) TableName() string {
	return table
}

func (db *store) Create(record *Record, session_length int) error {
	var err error
	// check for required attrs
	if record.UserId == 0 {
		return fmt.Errorf("User ID must not be blank")
	}

	// TODO: while unlikely, duplicate tokens would cause db layer error
	record.Token = util.RandSeq(20, util.LowerAlphaNumeric)
	record.ExpireTime = time.Now().UTC().AddDate(0, 0, session_length)

	query := fmt.Sprintf(`
        INSERT INTO %s
			(user_id, token, expire_time)
        VALUES
			(:user_id, :token, :expire_time) RETURNING id`, table)

	stmt, err := db.sqlx.PrepareNamed(query)
	err = stmt.QueryRowx(record).Scan(&record.Id)
	return err
}

func (db *store) GetByToken(token string) (Record, error) {
	record := Record{}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE token = ?", table))
	err := db.sqlx.Get(&record, query, token)
	return record, err
}

func (db *store) Authenticate(token string) (id int64, err error) {
	var record Record
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE token = ?", table))
	err = db.sqlx.Get(&record, query, token)
	if err != nil {
		return 0, err
	}

	if record.ExpireTime.UnixNano() < time.Now().UTC().UnixNano() {
		return 0, fmt.Errorf("Token expired.")
	}

	id = record.UserId
	return
}

func (db *store) Delete(token string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE token = ?", table)
	_, err := db.sqlx.Exec(db.sqlx.Rebind(query), token)
	return err
}