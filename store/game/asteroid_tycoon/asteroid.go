package asteroid_tycoon

import (
	"fmt"
	"time"

	"github.com/cbarraford/cryptocades-backend/util"
	"github.com/jmoiron/sqlx"
)

const asteroidsTable string = "g2_asteroids"

type Asteroid struct {
	Id          int64     `json:"id" db:"id"`
	Total       int       `json:"total" db:"total"`
	Remaining   int       `json:"remaining" db:"remaining"`
	Distance    int       `json:"distance" db:"distance"`
	ShipId      int64     `json:"ship_id" db:"ship_id"`
	SolarSystem int       `json:"solar_system" db:"solar_system"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
	UpdatedTime time.Time `json:"updated_time" db:"updated_time"`
}

func (db *store) CreateAsteroid(ast *Asteroid) error {
	var err error

	if ast.CreatedTime.IsZero() {
		ast.CreatedTime = time.Now().UTC()
	}

	if ast.Distance == 0 {
		ast.Distance = util.Random(10, 100)
	}

	if ast.Total == 0 {
		ast.Total = util.Random(100, 500)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s
			(total, remaining, distance, created_time)
		VALUES
			(:total, :remaining, :distance, :created_time) RETURNING id
	`, asteroidsTable)

	stmt, err := db.sqlx.PrepareNamed(query)
	if err != nil {
		return err
	}
	err = stmt.QueryRowx(ast).Scan(&ast.Id)
	return err
}

func (db *store) Mined(sessionId string, shares int, tx *sqlx.Tx) error {
	var err error

	query := db.sqlx.Rebind(fmt.Sprintf(`
		UPDATE %s AS ast SET
			remaining = remaining - ?,
			updated_time = now()
		FROM %s AS sessions
		WHERE sessions.session_id = ? AND sessions.ship_id = ast.ship_id
	`, asteroidsTable, ledgersTable))

	_, err = tx.Exec(query, shares*ResourceToShareRatio, sessionId)
	return err
}

func (db *store) OwnedAsteroid(shipId int64) (Asteroid, error) {
	asteroid := Asteroid{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE ship_id = ?", asteroidsTable)
	query = db.sqlx.Rebind(query)
	err := db.sqlx.Get(&asteroid, query, shipId)
	return asteroid, err
}

func (db *store) AvailableAsteroids() ([]Asteroid, error) {
	asteroids := []Asteroid{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE ship_id = 0 ORDER BY id DESC", asteroidsTable)
	err := db.sqlx.Select(&asteroids, query)
	return asteroids, err
}

func (db *store) AssignAsteroid(id, shipId int64) error {
	query := fmt.Sprintf(`
        UPDATE %s SET
            ship_id         = ?,
            updated_time    = now()
        WHERE id = ? AND ship_id = 0`, asteroidsTable)
	query = db.sqlx.Rebind(query)
	_, err := db.sqlx.Exec(query, shipId, id)
	return err
}

func (db *store) DestroyAsteroids() error {
	query := fmt.Sprintf("DELETE FROM %s WHERE remaining < 1", asteroidsTable)
	_, err := db.sqlx.Exec(db.sqlx.Rebind(query))
	return err
}
