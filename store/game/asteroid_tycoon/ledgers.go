package asteroid_tycoon

import (
	"fmt"
	"time"
)

const (
	ledgersTable         string = "g2_ledgers"
	sessionsTable        string = "g2_sessions"
	ResourceToShareRatio int    = 100
)

type Ledger struct {
	Id          int64     `json:"id" db:"id"`
	AccountId   int64     `json:"account_id" db:"account_id"`
	SessionId   string    `json:"session_id" db:"session_id"`
	Amount      int       `json:"amount" db:"amount"`
	Description string    `json:"Description" db:"description"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
}

func (db *store) Complete(sessionId string, shares int) error {
	var err error

	tx, err := db.sqlx.Beginx()

	query := db.sqlx.Rebind(fmt.Sprintf(`
		UPDATE %s AS ast SET
			remaining = remaining - ?
		FROM %s AS sessions
		WHERE sessions.session_id = ? AND sessions.ship_id = ast.ship_id
	`, asteroidsTable, ledgersTable))

	_, err = tx.Exec(query, shares*ResourceToShareRatio, sessionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
