package income

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
)

type Store interface {
	Create(record *Record) error
	Get(id int64) (Record, error)
	ListByUser(userId int64) ([]Record, error)
	UserIncome(userId int64) (int, error)
	UpdateScores() error
}

type store struct {
	Store
	sqlx  *sqlx.DB
	redis redis.Conn
}

func NewStore(db *sqlx.DB, redis redis.Conn) Store {
	return &store{sqlx: db, redis: redis}
}

const table string = "incomes"

type Record struct {
	Id          int64     `json:"id" db:"id"`
	UserId      int64     `json:"user_id" db:"user_id"`
	GameId      int64     `json:"game_id" db:"game_id"`
	SessionId   string    `json:"session_id" db:"session_id"`
	Amount      int       `json:"amount" db:"amount"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
}

func (db *store) TableName() string {
	return table
}

func (db *store) Create(record *Record) error {
	var err error
	// check for required attrs
	if record.SessionId == "" {
		return fmt.Errorf("Session id must not be blank")
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

	err = db.CreateWithinTransaction(record, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *store) CreateWithinTransaction(record *Record, tx *sqlx.Tx) error {
	query := db.sqlx.Rebind(fmt.Sprintf(`
        INSERT INTO %s
            (game_id, session_id, user_id, amount)
        VALUES
            (?, ?, ?, ?) ON CONFLICT (game_id, session_id, user_id) DO UPDATE SET amount = %s.amount + ?`, table, table))

	_, err := tx.Exec(query, record.GameId, record.SessionId, record.UserId, record.Amount, record.Amount)
	return err
}

func (db *store) Get(id int64) (Record, error) {
	record := Record{}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE id = ?", table))
	err := db.sqlx.Get(&record, query, id)
	return record, err
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

func (db *store) UserIncome(userId int64) (spent int, err error) {
	// If you update the query here, also update it in entry.go. This isn't
	// very DRY, but felt like it wasn't worth the import, getting this func to
	// work within a transaction, etc
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT COALESCE(SUM(amount),0) FROM %s WHERE user_id = ?", table))
	err = db.sqlx.Get(&spent, query, userId)
	return
}
