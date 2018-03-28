package asteroid_tycoon

import (
	"fmt"
	"time"
)

const accountsTable string = "g2_accounts"

type Account struct {
	Id          int64     `json:"id" db:"id"`
	UserId      int64     `json:"user_id" db:"user_id"`
	Credits     int       `json:"credits" db:"credits"`
	Resources   int       `json:"resources" db:"resources"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
	UpdatedTime time.Time `json:"updated_time" db:"updated_time"`
}

func (db *store) CreateAccount(acct *Account) error {
	var err error
	// check for required attrs
	if acct.UserId == 0 {
		return fmt.Errorf("Must belong to a user.")
	}

	if acct.CreatedTime.IsZero() {
		acct.CreatedTime = time.Now().UTC()
	}

	query := fmt.Sprintf(`
		INSERT INTO %s
			(user_id, created_time)
		VALUES
			(:user_id, :created_time) RETURNING id
	`, accountsTable)

	stmt, err := db.sqlx.PrepareNamed(query)
	err = stmt.QueryRowx(acct).Scan(&acct.Id)
	return err
}

func (db *store) GetAccountByUserId(userId int64) (acct Account, err error) {
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE user_id = ?", accountsTable))
	err = db.sqlx.Get(&acct, query, userId)
	return acct, err
}

func (db *store) UpdateAccount(acct *Account) error {
	// touch updated time
	acct.UpdatedTime = time.Now()

	query := fmt.Sprintf(`
        UPDATE %s SET  
            credits         = :credits,
            resources       = :resources,
            updated_time    = :updated_time
        WHERE id = :id`, accountsTable)
	_, err := db.sqlx.NamedExec(query, acct)
	return err
}

func (db *store) DeleteAccount(id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", accountsTable)
	_, err := db.sqlx.Exec(db.sqlx.Rebind(query), id)
	return err
}
