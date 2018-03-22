package boost

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	Create(record *Record) error
	Get(id int64) (Record, error)
	ListByUser(userId int64) ([]Record, error)
	Assign(id, income_id int64) error
}

type store struct {
	Store
	sqlx *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return &store{sqlx: db}
}

const table string = "boosts"

type Record struct {
	Id          int64     `json:"id" db:"id"`
	UserId      int64     `json:"user_id" db:"user_id"`
	IncomeId    int64     `json:"income_id" db:"income_id"`
	Multiplier  int       `json:"multiplier" db:"multiplier"`
	UpdatedTime time.Time `json:"updated_time" db:"updated_time"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
}

func (db *store) TableName() string {
	return table
}

func (db *store) Create(record *Record) error {
	var err error
	// check for required attrs
	if record.UserId == 0 {
		return fmt.Errorf("User id must not be blank")
	}

	now := time.Now()
	if record.CreatedTime.IsZero() {
		record.CreatedTime = now
	}
	if record.UpdatedTime.IsZero() {
		record.UpdatedTime = now
	}

	// we're purposely ignore multipler and income_id for security purposes. We
	// don't want a boost to be created an applied to an income in the same
	// transaction. We also don't want a hacker to manipulate our multipler.
	query := fmt.Sprintf(`
        INSERT INTO %s
			(user_id, created_time, updated_time)
        VALUES
			(:user_id, :created_time, :updated_time) RETURNING id`, table)

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

func (db *store) List() (records []Record, err error) {
	query := fmt.Sprintf("SELECT * FROM %s", table)
	err = db.sqlx.Select(&records, query)
	return
}

func (db *store) ListByUser(id int64) ([]Record, error) {
	records := []Record{}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE user_id = ?", table))
	err := db.sqlx.Select(&records, query, id)
	return records, err
}

func (db *store) Assign(id, income_id int64) error {
	record, err := db.Get(id)
	if err != nil {
		return err
	}

	// can't assign a boost that is already assigned
	if record.IncomeId > 0 {
		return fmt.Errorf("This boost is already assigned to a previous game session.")
	}

	var incUserId int64
	query := db.sqlx.Rebind("SELECT user_id FROM incomes WHERE id = ? AND game_id > 0")
	err = db.sqlx.Get(&incUserId, query, income_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Income id not found.")
		}
		return err
	}

	if incUserId != record.UserId {
		return fmt.Errorf("This boost and income session is not owned by the same user.")
	}

	// do the assignment
	record.IncomeId = income_id

	// touch updated time
	record.UpdatedTime = time.Now()

	query = fmt.Sprintf(`
        UPDATE %s SET
            income_id		= :income_id,
            updated_time	= :updated_time
        WHERE id = :id`, table)
	_, err = db.sqlx.NamedExec(query, record)
	return err
}
