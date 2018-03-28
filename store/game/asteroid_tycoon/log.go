package asteroid_tycoon

import (
	"fmt"
	"time"
)

const logsTable = "g2_logs"

type Log struct {
	Id          int64     `json:"id" db:"id"`
	ShipId      int64     `json:"ship_id" db:"ship_id"`
	Level       int       `json:"level" db:"level"`
	Log         string    `json:"log" db:"log"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
}

func (db *store) CreateLog(line *Log) error {
	var err error
	// check for required attrs
	if line.ShipId == 0 {
		return fmt.Errorf("Must belong to a ship.")
	}
	if line.Log == "" {
		return fmt.Errorf("Cannot create a blank log line.")
	}

	if line.CreatedTime.IsZero() {
		line.CreatedTime = time.Now().UTC()
	}

	query := fmt.Sprintf(`
		INSERT INTO %s
			(ship_id, level, log, created_time)
		VALUES
			(:ship_id, :level, :log, :created_time) RETURNING id
	`, logsTable)

	stmt, err := db.sqlx.PrepareNamed(query)
	if err != nil {
		return err
	}
	err = stmt.QueryRowx(line).Scan(&line.Id)
	return err
}

func (db *store) GetShipLogs(shipId int64) ([]Log, error) {
	lines := []Log{}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE ship_id = ?", logsTable))
	err := db.sqlx.Select(&lines, query, shipId)
	return lines, err
}

func (db *store) DeleteOldLogs() error {
	query := fmt.Sprintf("DELETE FROM %s WHERE created_time < (now() - '7 day'::INTERVAL)", accountsTable)
	_, err := db.sqlx.Exec(db.sqlx.Rebind(query))
	return err
}
