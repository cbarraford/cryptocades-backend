package user

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	Create(record *Record) error
	Get(id int64) (Record, error)
	GetByUsername(username string) (Record, error)
	GetByBTCAddress(btc_address string) (Record, error)
	Update(record *Record) error
	List() ([]Record, error)
	Authenticate(username, password string) (Record, error)
}

type store struct {
	Store
	sqlx *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return &store{sqlx: db}
}

const table string = "users"

type Record struct {
	Id          int64     `json:"id" db:"id"`
	BTCAddr     string    `json:"btc_address" db:"btc_address"`
	Username    string    `json:"username" db:"username"`
	Password    string    `json:"-" db:"password"`
	Email       string    `json:"email" db:"email"`
	MinedHashes int       `json:"mined_hashes" db:"mined_hashes"`
	BonusHashes int       `json:"bonus_hashes" db:"bonus_hashes"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
	UpdatedTime time.Time `json:"updated_time" db:"updated_time"`
}

func (db *store) TableName() string {
	return table
}

func (db *store) Create(record *Record) error {
	var err error
	// check for required attrs
	if record.Username == "" {
		return fmt.Errorf("Username must not be blank")
	}
	if record.Email == "" {
		return fmt.Errorf("Email must not be blank")
	}
	if record.BTCAddr == "" {
		return fmt.Errorf("Bitcoin Address must not be blank")
	}
	if record.Password == "" {
		return fmt.Errorf("Password must not be blank")
	}

	// always store password hashed and salted
	record.Password, err = HashPassword(record.Password)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
        INSERT INTO %s
			(username, password, btc_address, email)
        VALUES
			(:username, :password, :btc_address, :email) RETURNING id`, table)

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

func (db *store) GetByUsername(username string) (record Record, err error) {
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE username = ?", table))
	err = db.sqlx.Get(&record, query, username)
	return
}

func (db *store) GetByBTCAddress(btc string) (record Record, err error) {
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE btc_address = ?", table))
	err = db.sqlx.Get(&record, query, btc)
	return
}

// TODO: change password function

func (db *store) Update(record *Record) error {
	// touch updated time
	record.UpdatedTime = time.Now()

	query := fmt.Sprintf(`
        UPDATE %s SET
            username        = :username,
            email           = :email,
            updated_time    = :updated_time,
			mined_hashes	= :mined_hashes,
			bonus_hashes	= :bonus_hashes,
			btc_address		= :btc_address
        WHERE id = :id`, table)
	_, err := db.sqlx.NamedExec(query, record)
	return err
}

func (db *store) List() (records []Record, err error) {
	query := fmt.Sprintf("SELECT * FROM %s", table)
	err = db.sqlx.Select(&records, query)
	return
}

func (db *store) Authenticate(username, password string) (record Record, err error) {
	incorrect := fmt.Errorf("Incorrect username or password")
	record, err = db.GetByUsername(username)
	if err != nil {
		return Record{}, incorrect
	}

	if !CheckPasswordHash(password, record.Password) {
		return Record{}, incorrect
	}

	return
}
