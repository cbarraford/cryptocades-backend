package income

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"

	"github.com/cbarraford/cryptocades-backend/store/game/asteroid_tycoon"
)

type Store interface {
	Create(record *Record) error
	Get(id int64) (Record, error)
	ListByUser(userId int64) ([]Record, error)
	UserIncome(userId int64) (int, error)
	UserIncomeRank(userId int64) (int, error)
	UpdateScores(tyGame asteroid_tycoon.Store) error
	CountBonuses(userId int64, prefix string) (int, error)
}

const rankDivider = 100

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
	Id            int64     `json:"id" db:"id"`
	UserId        int64     `json:"user_id" db:"user_id"`
	GameId        int64     `json:"game_id" db:"game_id"`
	SessionId     string    `json:"session_id" db:"session_id"`
	Amount        int       `json:"amount" db:"amount"`
	PartialAmount int       `json:"partial_amount" db:"partial_amount"`
	UpdatedTime   time.Time `json:"updated_time" db:"updated_time"`
	CreatedTime   time.Time `json:"created_time" db:"created_time"`
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

	if record.CreatedTime.IsZero() {
		record.CreatedTime = time.Now()
	}
	if record.UpdatedTime.IsZero() {
		record.UpdatedTime = time.Now()
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
	// TODO: we are hard coding the spread for Tallest Tower. Will need to
	// store this value in each individual game and send the game's spread via
	// DI
	spread := 20
	query := db.sqlx.Rebind(fmt.Sprintf(`
        INSERT INTO %s
            (game_id, session_id, user_id, amount, partial_amount)
        VALUES
            (?, ?, ?, ?, ?) ON CONFLICT (game_id, session_id, user_id) DO UPDATE SET amount = %s.amount + ? + FLOOR((%s.partial_amount + ?)/?), partial_amount = MOD(%s.partial_amount + ?,?), updated_time = now()`, table, table, table, table))

	_, err := tx.Exec(
		query,
		record.GameId,
		record.SessionId,
		record.UserId,
		record.Amount+(record.PartialAmount/spread),
		record.PartialAmount%spread,
		record.Amount,
		record.PartialAmount,
		spread,
		record.PartialAmount,
		spread,
	)
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

func (db *store) ListByUser(userId int64) (records []Record, err error) {
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE user_id = ?", table))
	err = db.sqlx.Select(&records, query, userId)
	return
}

func (db *store) CountBonuses(id int64, prefix string) (i int, err error) {
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE session_id LIKE ?", table))
	err = db.sqlx.Get(&i, query, prefix+"%")
	return
}

func (db *store) UserIncome(userId int64) (earned int, err error) {
	// If you update the query here, also update it in entry.go. This isn't
	// very DRY, but felt like it wasn't worth the import, getting this func to
	// work within a transaction, etc
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT COALESCE(SUM(%s.amount * COALESCE(boosts.multiplier,1)),0) FROM %s LEFT OUTER JOIN boosts ON boosts.income_id = %s.id WHERE %s.user_id = ?", table, table, table, table))
	err = db.sqlx.Get(&earned, query, userId)
	return
}

func (db *store) UserIncomeRank(userId int64) (rank int, err error) {
	query := db.sqlx.Rebind(fmt.Sprintf("WITH totals AS (SELECT user_id, COALESCE(SUM(amount),0) AS total FROM %s GROUP BY user_id), ranks AS (SELECT user_id, total, ntile(%d) OVER (ORDER BY total) AS rank from totals) SELECT rank from ranks WHERE user_id = ?;", table, rankDivider))
	_ = db.sqlx.Get(&rank, query, userId)
	rank = int((float64(rank) / float64(rankDivider)) * 100)
	return
}
