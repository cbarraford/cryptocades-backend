package asteroid_tycoon

import (
	"fmt"
	"time"

	"github.com/cbarraford/cryptocades-backend/util"
	"github.com/lib/pq"
)

const accountsTable string = "g2_accounts"

type Account struct {
	Id          int64     `json:"id" db:"id"`
	UserId      int64     `json:"user_id" db:"user_id"`
	Credits     int       `json:"credits" db:"credits"`
	Resources   int       `json:"resources"`
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
	if err != nil {
		return acct, err
	}

	acct.Resources, err = db.ResourceBalance(acct.Id)
	return acct, err
}

func (db *store) UpdateAccount(acct *Account) error {
	query := fmt.Sprintf(`
        UPDATE %s SET  
            credits         = :credits,
            updated_time    = now()
        WHERE id = :id`, accountsTable)
	_, err := db.sqlx.NamedExec(query, acct)
	return err
}

func (db *store) AddAccountCredits(id int64, amount int) error {
	query := db.sqlx.Rebind(fmt.Sprintf(`
        UPDATE %s SET  
            credits       = credits + ?,
            updated_time    = now()
        WHERE id = ?`, accountsTable))
	_, err := db.sqlx.Exec(query, amount, id)
	return err
}

func (db *store) SubtractAccountCredits(id int64, amount int) error {
	query := db.sqlx.Rebind(fmt.Sprintf(`
        UPDATE %s SET  
            credits       = credits - ?,
            updated_time    = now()
        WHERE id = ?`, accountsTable))
	_, err := db.sqlx.Exec(query, amount, id)
	return err
}

func (db *store) DeleteAccount(id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", accountsTable)
	_, err := db.sqlx.Exec(db.sqlx.Rebind(query), id)
	return err
}

func (db *store) TradeForPlays(accountId int64, amount int) error {
	creditNeeded := CreditsForPlays * amount
	tx, err := db.sqlx.Beginx()

	query := db.sqlx.Rebind(fmt.Sprintf(`
        UPDATE %s SET
			credits		= credits - ?,
            updated_time    = now()
        WHERE id = ?`, accountsTable))
	_, err = tx.Exec(query, creditNeeded, accountId)
	if serr, ok := err.(*pq.Error); ok {
		// https://www.postgresql.org/docs/9.3/static/errcodes-appendix.html
		if serr.Code.Name() == "check_violation" {
			tx.Rollback()
			return fmt.Errorf("Insufficient funds.")
		}
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	var userId int64

	query = db.sqlx.Rebind(fmt.Sprintf("SELECT user_id FROM %s WHERE id = ?", accountsTable))
	err = tx.Get(&userId, query, accountId)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = db.sqlx.Rebind(`
        INSERT INTO incomes
            (game_id, session_id, user_id, amount, partial_amount)
        VALUES
            (?, ?, ?, ?, 0) ON CONFLICT (game_id, session_id, user_id) DO UPDATE SET amount = incomes.amount + ?, updated_time = now()`)

	_, err = tx.Exec(
		query,
		2,
		util.RandSeq(16, util.LowerAlphaNumeric),
		userId,
		amount,
		amount,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *store) TradeForCredits(accountId int64, amount int) error {
	var query string

	tx, err := db.sqlx.Beginx()

	var balance int
	query = db.sqlx.Rebind(fmt.Sprintf("SELECT COALESCE(SUM(amount),0) FROM %s WHERE account_id = ?", ledgersTable))
	err = db.sqlx.Get(&balance, query, accountId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if balance < amount*ResourcesForCredits {
		tx.Rollback()
		return fmt.Errorf("Insufficient funds.")
	}

	query = db.sqlx.Rebind(fmt.Sprintf(`
		INSERT INTO %s AS ledger
			(account_id, session_id, amount, description)
		VALUES
			(?, ?, ?, ?)`, ledgersTable))
	_, err = tx.Exec(
		query,
		accountId,
		util.RandSeq(12, util.LowerAlphaNumeric),
		-amount*ResourcesForCredits,
		"Trade for credits",
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = db.sqlx.Rebind(fmt.Sprintf(`
        UPDATE %s SET
			credits		= credits + ?,
            updated_time    = now()
        WHERE id = ?`, accountsTable))
	_, err = tx.Exec(query, amount, accountId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
