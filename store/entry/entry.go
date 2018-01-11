package entry

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	Create(record *Record) error
	Get(id int64) (Record, error)
	GetOdds(jackpotId, userId int64) (Odds, error)
	List() ([]Record, error)
	ListByUser(id int64) ([]Record, error)
	UserSpent(userId int64) (int, error)
}

type store struct {
	Store
	sqlx *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return &store{sqlx: db}
}

const table string = "entries"

type Record struct {
	Id          int64     `json:"id" db:"id"`
	JackpotId   int64     `json:"jackpot_id" db:"jackpot_id"`
	UserId      int64     `json:"user_id" db:"user_id"`
	Amount      int       `json:"amount" db:"amount"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
}

type Odds struct {
	JackpotId int64 `json:"jackpot_id"`
	Total     int64 `json:"total"`
	Entries   int64 `json:"entries"`
}

func (db *store) TableName() string {
	return table
}

func (db *store) Create(record *Record) error {
	var err error
	// check for required attrs
	if record.JackpotId == 0 {
		return fmt.Errorf("Jackpot id must not be blank")
	}
	if record.UserId == 0 {
		return fmt.Errorf("User id must not be blank")
	}
	if record.Amount == 0 {
		return fmt.Errorf("Amount must be more than 0")
	}

	var count int64
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE user_id = ? AND jackpot_id = ?", table))
	err = db.sqlx.Get(&count, query, record.UserId, record.JackpotId)
	if err != nil {
		return err
	}

	if count > 0 {
		query := db.sqlx.Rebind(
			fmt.Sprintf("UPDATE %s SET amount = amount + ? WHERE user_id = ? AND jackpot_id = ?", table))
		_, err = db.sqlx.Exec(query, record.Amount, record.UserId, record.JackpotId)
	} else {

		query := fmt.Sprintf(`
        INSERT INTO %s
			(jackpot_id, user_id, amount)
        VALUES
			(:jackpot_id, :user_id, :amount) RETURNING id`, table)

		stmt, err := db.sqlx.PrepareNamed(query)
		if err != nil {
			return err
		}
		err = stmt.QueryRowx(record).Scan(&record.Id)
	}
	return err
}

func (db *store) Get(id int64) (Record, error) {
	record := Record{}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE id = ?", table))
	err := db.sqlx.Get(&record, query, id)
	return record, err
}

func (db *store) GetOdds(jackpotId, userId int64) (odd Odds, err error) {
	odd.JackpotId = jackpotId
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT COALESCE(SUM(amount),0) FROM %s WHERE jackpot_id = ?", table))
	err = db.sqlx.Get(&odd.Total, query, jackpotId)
	if err != nil {
		return
	}

	query = db.sqlx.Rebind(fmt.Sprintf("SELECT COALESCE(SUM(amount),0) FROM %s WHERE jackpot_id = ? AND user_id = ?", table))
	err = db.sqlx.Get(&odd.Entries, query, jackpotId, userId)
	return
}

func (db *store) List() (records []Record, err error) {
	query := fmt.Sprintf("SELECT * FROM %s", table)
	err = db.sqlx.Select(&records, query)
	return
}

func (db *store) ListByUser(id int64) (records []Record, err error) {
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE user_id = ?", table))
	err = db.sqlx.Select(&records, query, id)
	return
}

func (db *store) UserSpent(userId int64) (spent int, err error) {
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT COALESCE(SUM(amount),0) FROM %s WHERE user_id = ?", table))
	err = db.sqlx.Get(&spent, query, userId)
	return
}
