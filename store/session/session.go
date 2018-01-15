package session

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/CBarraford/lotto/util"
)

const (
	escalatedTime = 5
)

type Store interface {
	Create(record *Record, i int) error
	GetByToken(token string) (Record, error)
	Authenticate(token string) (int64, bool, error)
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
	Id            int64     `json:"id" db:"id"`
	UserId        int64     `json:"user_id" db:"user_id"`
	Token         string    `json:"token" db:"token"`
	CreatedTime   time.Time `json:"created_time" db:"created_time"`
	EscalatedTime time.Time `json:"escalated_time"`
	ExpireTime    time.Time `json:"expire_time" db:"expire_time"`
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

	if record.CreatedTime.IsZero() {
		record.CreatedTime = time.Now().UTC()
	}
	record.EscalatedTime = record.CreatedTime.Add(escalatedTime * time.Minute)

	// TODO: while unlikely, duplicate tokens would cause db layer error
	record.Token = util.RandSeq(20, util.LowerAlphaNumeric)
	record.ExpireTime = time.Now().UTC().AddDate(0, 0, session_length)

	query := fmt.Sprintf(`
        INSERT INTO %s
			(user_id, token, expire_time, created_time)
        VALUES
			(:user_id, :token, :expire_time, :created_time) RETURNING id`, table)

	stmt, err := db.sqlx.PrepareNamed(query)
	err = stmt.QueryRowx(record).Scan(&record.Id)
	return err
}

func (db *store) GetByToken(token string) (Record, error) {
	record := Record{}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE token = ?", table))
	err := db.sqlx.Get(&record, query, token)
	record.EscalatedTime = record.CreatedTime.Add(escalatedTime * time.Minute)
	return record, err
}

func (db *store) Authenticate(token string) (id int64, escalatedPriv bool, err error) {
	var record Record
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE token = ?", table))
	err = db.sqlx.Get(&record, query, token)
	if err != nil {
		return 0, false, err
	}

	if record.ExpireTime.UnixNano() < time.Now().UTC().UnixNano() {
		return 0, false, fmt.Errorf("Token expired.")
	}

	if record.CreatedTime.Add(escalatedTime*time.Minute).UnixNano() >= time.Now().UnixNano() {
		escalatedPriv = true
	}

	id = record.UserId
	return
}

func (db *store) Delete(token string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE token = ?", table)
	_, err := db.sqlx.Exec(db.sqlx.Rebind(query), token)
	return err
}
