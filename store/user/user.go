package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/cbarraford/cryptocades-backend/util/gravatar"
)

type Store interface {
	Create(record *Record) error
	Get(id int64) (Record, error)
	GetByUsername(username string) (Record, error)
	GetByEmail(email string) (Record, error)
	GetByBTCAddress(btc_address string) (Record, error)
	GetByReferralCode(code string) (Record, error)
	GetByFacebookId(id string) (Record, error)
	Update(record *Record) error
	List() ([]Record, error)
	MarkAsConfirmed(record *Record) error
	Authenticate(username, password string) (Record, error)
	PasswordSet(record *Record) error
	Delete(id int64) error
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
	Id           int64     `json:"id" db:"id"`
	BTCAddr      string    `json:"btc_address" db:"btc_address"`
	Username     string    `json:"username" db:"username"`
	Password     string    `json:"-" db:"password"`
	Email        string    `json:"email" db:"email"`
	FacebookId   string    `json:"-" db:"fb_id"`
	Avatar       string    `json:"avatar" db:"-"`
	Confirmed    bool      `json:"confirmed" db:"confirmed"`
	ReferralCode string    `json:"referral_code" db:"referral_code"`
	Admin        bool      `json:"-" db:"admin"`
	CreatedTime  time.Time `json:"created_time" db:"created_time"`
	UpdatedTime  time.Time `json:"updated_time" db:"updated_time"`
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
	if record.Password == "" && record.FacebookId == "" {
		return fmt.Errorf("Password must not be blank")
	}

	// all users are auto-confirmed in development environment
	if os.Getenv("ENVIRONMENT") == "development" {
		record.Confirmed = true
	}

	// always store password hashed and salted
	record.Password, err = HashPassword(record.Password)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
        INSERT INTO %s
			(username, password, btc_address, email, fb_id, confirmed)
        VALUES
			(:username, :password, :btc_address, :email, :fb_id, :confirmed) RETURNING id`, table)

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

func (db *store) GetByEmail(email string) (record Record, err error) {
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE email = ?", table))
	err = db.sqlx.Get(&record, query, email)
	return
}

func (db *store) GetByReferralCode(code string) (record Record, err error) {
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE referral_code = ?", table))
	err = db.sqlx.Get(&record, query, code)
	return
}

func (db *store) GetByFacebookId(id string) (record Record, err error) {
	if id == "" {
		return record, fmt.Errorf("Facebook id not found")
	}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE fb_id = ?", table))
	err = db.sqlx.Get(&record, query, id)
	return
}

func (db *store) Update(record *Record) error {
	// touch updated time
	record.UpdatedTime = time.Now()

	query := fmt.Sprintf(`
        UPDATE %s SET
            username        = :username,
            email           = :email,
            updated_time    = :updated_time,
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

func (db *store) PasswordSet(record *Record) error {
	var err error

	// touch updated time
	record.UpdatedTime = time.Now()
	// since password reset utilizes email address, we are inherently
	// confirming the account.
	record.Confirmed = true

	// always store password hashed and salted
	record.Password, err = HashPassword(record.Password)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
        UPDATE %s SET
            confirmed       = :confirmed,
            updated_time    = :updated_time,
			password		= :password
        WHERE id = :id`, table)
	_, err = db.sqlx.NamedExec(query, record)
	return err
}

// Marks account as confirmed as well updates email address (if confirming new
// email address)
func (db *store) MarkAsConfirmed(record *Record) error {
	// touch updated time
	record.UpdatedTime = time.Now()
	record.Confirmed = true

	query := fmt.Sprintf(`
        UPDATE %s SET
			email			= :email,
            confirmed       = :confirmed,
            updated_time    = :updated_time
        WHERE id = :id`, table)
	_, err := db.sqlx.NamedExec(query, record)
	return err
}

func (db *store) Authenticate(username, password string) (record Record, err error) {
	incorrect := fmt.Errorf("Incorrect username or password")
	record, err = db.GetByUsername(username)

	// if password is blank, always fail
	if password == "" {
		return Record{}, incorrect
	}

	if err != nil || !record.Confirmed {
		return Record{}, incorrect
	}

	if !CheckPasswordHash(password, record.Password) {
		return Record{}, incorrect
	}

	return
}

func (db *store) Delete(id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", table)
	_, err := db.sqlx.Exec(db.sqlx.Rebind(query), id)
	return err
}

func (r Record) MarshalJSON() ([]byte, error) {
	type Alias Record
	t := struct {
		Avatar string `json:"avatar"`
		Alias
	}{
		Avatar: gravatar.Avatar(r.Email, 256),
		Alias:  (Alias)(r),
	}

	b := bytes.NewBuffer([]byte(``))
	e := json.NewEncoder(b)
	e.SetEscapeHTML(false)
	err := e.Encode(t)
	return b.Bytes(), err
}
