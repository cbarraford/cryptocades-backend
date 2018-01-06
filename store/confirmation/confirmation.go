package confirmation

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	Create(record *Record) error
	Get(id int64) (Record, error)
	GetByCode(code string) (Record, error)
	Delete(id int64) error
}

type store struct {
	Store
	sqlx *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return &store{sqlx: db}
}

const table string = "confirmations"

type Record struct {
	Id          int64     `json:"id" db:"id"`
	Code        string    `json:"code" db:"code"`
	Email       string    `json:"email" db:"email"`
	UserId      int64     `json:"user_id" db:"user_id"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
}

func (db *store) TableName() string {
	return table
}

func (db *store) Create(record *Record) error {
	var err error
	// check for required attrs
	if record.Code == "" {
		return fmt.Errorf("Code must not be blank")
	}
	if record.Email == "" {
		return fmt.Errorf("Email must not be blank")
	}
	if record.UserId == 0 {
		return fmt.Errorf("UserId must not be zero")
	}

	query := fmt.Sprintf(`
        INSERT INTO %s
			(code, email, user_id)
        VALUES
			(:code, :email, :user_id) RETURNING id`, table)

	stmt, err := db.sqlx.PrepareNamed(query)
	err = stmt.QueryRowx(record).Scan(&record.Id)
	return err
}

func (db *store) Get(id int64) (Record, error) {
	record := Record{}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE id = ?", table))
	err := db.sqlx.Get(&record, query, id)
	return record, err
}

func (db *store) GetByCode(code string) (record Record, err error) {
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE code = ?", table))
	err = db.sqlx.Get(&record, query, code)
	return
}

func (db *store) Delete(id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", table)
	_, err := db.sqlx.Exec(db.sqlx.Rebind(query), id)
	return err
}
