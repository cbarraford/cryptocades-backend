package entry

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	Create(record *Record) error
	Get(id int64) (Record, error)
	GetOdds(jackpotId, userId int64) (Odds, error)
	List() ([]Record, error)
	ListByUser(id int64) ([]Record, error)
	ListByJackpot(id int64) ([]Record, error)
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

	tx, err := db.sqlx.Beginx()
	if err != nil {
		return err
	}

	query := db.sqlx.Rebind(fmt.Sprintf(`
        INSERT INTO %s
			(jackpot_id, user_id, amount)
        VALUES
			(?, ?, ?) ON CONFLICT (user_id, jackpot_id) DO UPDATE SET amount = %s.amount + ?`, table, table))

	_, err = tx.Exec(query, record.JackpotId, record.UserId, record.Amount, record.Amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	var totalIncome int
	// this query should stay in sync with income.go:UserIncome
	query = db.sqlx.Rebind("SELECT COALESCE(SUM(incomes.amount * COALESCE(boosts.multiplier,1)),0) FROM incomes LEFT OUTER JOIN boosts ON boosts.income_id = incomes.id WHERE incomes.user_id = ?")
	err = tx.Get(&totalIncome, query, record.UserId)
	if err != nil {
		tx.Rollback()
		return err
	}

	var spent int
	query = db.sqlx.Rebind(fmt.Sprintf("SELECT COALESCE(SUM(amount),0) FROM %s WHERE user_id = ?", table))
	err = tx.Get(&spent, query, record.UserId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// check if we're overspent
	log.Printf("Total: %d | Spent: %d", totalIncome, spent)
	if (totalIncome - spent) < 0 {
		tx.Rollback()
		return fmt.Errorf("Insufficient funds.")
	}

	return tx.Commit()
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

func (db *store) ListByJackpot(id int64) (records []Record, err error) {
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE jackpot_id = ?", table))
	err = db.sqlx.Select(&records, query, id)
	return
}

func (db *store) UserSpent(userId int64) (spent int, err error) {
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT COALESCE(SUM(amount),0) FROM %s WHERE user_id = ?", table))
	err = db.sqlx.Get(&spent, query, userId)
	return
}
