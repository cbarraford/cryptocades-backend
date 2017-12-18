package jackpot

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	Create(record *Record) error
	Get(id int64) (Record, error)
	Update(record *Record) error
	List() ([]Record, error)
}

type store struct {
	Store
	sqlx *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return &store{sqlx: db}
}

const table string = "jackpots"

type Record struct {
	Id          int64     `json:"id" db:"id"`
	Jackpot     int       `json:"jackpot" db:"jackpot"`
	WinnerId    int64     `json:"winner_id" db:"winner_id"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
	UpdatedTime time.Time `json:"updated_time" db:"updated_time"`
}

func (db *store) TableName() string {
	return table
}

func (db *store) Create(record *Record) error {
	var err error
	// check for required attrs
	if record.EndTime.IsZero() {
		return fmt.Errorf("End time must not be blank")
	}

	query := fmt.Sprintf(`
        INSERT INTO %s
			(jackpot, end_time)
        VALUES
			(:jackpot, :end_time) RETURNING id`, table)

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

func (db *store) Update(record *Record) error {
	// touch updated time
	record.UpdatedTime = time.Now()

	query := fmt.Sprintf(`
        UPDATE %s SET
            jackpot			= :jackpot,
            updated_time	= :updated_time,
			winner_id		= :winner_id
        WHERE id = :id`, table)
	_, err := db.sqlx.NamedExec(query, record)
	return err
}

func (db *store) List() (records []Record, err error) {
	query := fmt.Sprintf("SELECT * FROM %s", table)
	err = db.sqlx.Select(&records, query)
	return
}
