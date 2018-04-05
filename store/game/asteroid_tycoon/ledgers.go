package asteroid_tycoon

import (
	"fmt"
	"time"

	"github.com/cbarraford/cryptocades-backend/util"
)

const (
	ledgersTable         string = "g2_ledgers"
	sessionsTable        string = "g2_sessions"
	ResourceToShareRatio int    = 100
	ResourcesForCredits  int    = 1000
)

type Ledger struct {
	Id          int64     `json:"id" db:"id"`
	AccountId   int64     `json:"account_id" db:"account_id"`
	SessionId   string    `json:"session_id" db:"session_id"`
	Amount      int       `json:"amount" db:"amount"`
	Description string    `json:"Description" db:"description"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
}

func (db *store) CompletedAsteroid(ast Asteroid) error {
	var err error
	description := "Asteroid Mined"
	amount := ast.Total
	if ast.Remaining > 0 {
		amount = amount - ast.Remaining
	}

	status := db.GetStatus(ast)
	if status.Status != "Docked" {
		return fmt.Errorf("Cannot collect resource while the ship is not docked")
	}

	tx, err := db.sqlx.Beginx()
	if err != nil {
		return err
	}

	ship, err := db.GetShip(ast.ShipId)
	if err != nil {
		tx.Rollback()
		return err
	}

	query := db.sqlx.Rebind(fmt.Sprintf(`
		INSERT INTO %s
		   (account_id, session_id, amount, description)
	    VALUES
		   (?, ?, ?, ?)
		ON CONFLICT (account_id, session_id) 
		DO UPDATE SET amount = %s.amount + ?
	`, ledgersTable, ledgersTable))
	_, err = tx.Exec(query, ship.AccountId, ast.SessionId, amount, description, amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = db.sqlx.Rebind(
		fmt.Sprintf("DELETE FROM %s WHERE id = ?", asteroidsTable),
	)
	_, err = tx.Exec(query, ast.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *store) AddEntry(accountId int64, amount int, description string) error {
	var err error
	query := db.sqlx.Rebind(fmt.Sprintf(`
		INSERT INTO %s AS ledger
			(account_id, session_id, amount, description)
		VALUES
			(?, ?, ?, ?)`, ledgersTable))
	_, err = db.sqlx.Exec(
		query,
		accountId,
		util.RandSeq(12, util.LowerAlphaNumeric),
		amount,
		description,
	)
	return err
}

func (db *store) ListByAccountId(accountId int64) ([]Ledger, error) {
	ledgers := []Ledger{}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE account_id = ?", ledgersTable))
	err := db.sqlx.Select(&ledgers, query, accountId)
	return ledgers, err
}

func (db *store) ResourceBalance(accountId int64) (balance int, err error) {
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT COALESCE(SUM(amount),0) FROM %s WHERE account_id = ?", ledgersTable))
	err = db.sqlx.Get(&balance, query, accountId)
	return
}
