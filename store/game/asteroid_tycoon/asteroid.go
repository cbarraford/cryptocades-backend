package asteroid_tycoon

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/cbarraford/cryptocades-backend/util"
	"github.com/jmoiron/sqlx"
)

const (
	asteroidsTable string = "g2_asteroids"
)

type Asteroid struct {
	Id          int64     `json:"id" db:"id"`
	Total       int       `json:"total" db:"total"`
	Remaining   int       `json:"remaining" db:"remaining"`
	Distance    int       `json:"distance" db:"distance"`
	ShipId      int64     `json:"ship_id" db:"ship_id"`
	ShipSpeed   int       `json:"ship_speed" db:"ship_speed"`
	SolarSystem int       `json:"solar_system" db:"solar_system"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
	UpdatedTime time.Time `json:"updated_time" db:"updated_time"`
}

func (db *store) CreateAsteroid(ast *Asteroid) error {
	var err error

	if ast.CreatedTime.IsZero() {
		ast.CreatedTime = time.Now().UTC()
	}
	if ast.UpdatedTime.IsZero() {
		ast.UpdatedTime = time.Now().UTC()
	}

	if ast.Distance == 0 {
		ast.Distance = util.Random(MinDistance, MaxDistance)
	}

	if ast.Total == 0 {
		ast.Total = util.Random(MinTotal, MaxTotal)
	}

	if ast.Remaining == 0 {
		ast.Remaining = ast.Total
	}

	query := fmt.Sprintf(`
		INSERT INTO %s
			(total, remaining, distance, updated_time, created_time)
		VALUES
			(:total, :remaining, :distance, :updated_time, :created_time) RETURNING id
	`, asteroidsTable)

	stmt, err := db.sqlx.PrepareNamed(query)
	if err != nil {
		return err
	}
	err = stmt.QueryRowx(ast).Scan(&ast.Id)
	return err
}

func (db *store) Mined(sessionId string, shares int, userId int64, tx *sqlx.Tx) error {
	var err error

	var ship Ship
	query := db.sqlx.Rebind(fmt.Sprintf(`
		SELECT * FROM %s WHERE session_id = ?
	`, shipsTable))
	err = tx.Get(&ship, query, sessionId)
	if err != nil {
		return err
	}
	err = db.ExpandShip(&ship)
	if err != nil {
		return err
	}

	// get the asteroid so we know the max damage
	var asteroid Asteroid
	query = db.sqlx.Rebind(fmt.Sprintf(`
		SELECT * FROM %s WHERE ship_id = ?
	`, asteroidsTable))
	err = tx.Get(&asteroid, query, ship.Id)

	shipStatus := db.GetStatus(ship, asteroid)
	var elapsedTime float64
	if shipStatus.TravelTime == 0 {
		elapsedTime = float64(time.Now().Unix() - asteroid.UpdatedTime.Unix())
	} else {
		elapsedTime = math.Min(float64(time.Now().Unix()-asteroid.UpdatedTime.Unix()), float64(shipStatus.TravelTime*2))
	}
	if err != nil {
		if err == sql.ErrNoRows {
			incr := int64(math.Min(float64(ship.Hull-ship.Health), elapsedTime*float64(ship.Repair)))
			if incr > 0 {
				query = db.sqlx.Rebind(fmt.Sprintf(`
				UPDATE %s SET
					health = health + ?
				WHERE id = ?
				`, shipsTable))
				_, err = tx.Exec(query, incr, ship.Id)
			}
			return err
		}
		return err
	}

	// ensure we can't mine when our health zero or below
	if ship.Health <= 0 {
		err = db.CompletedAsteroid(asteroid)
		if err != nil {
			log.Printf("Error completing asteroid: %+v", err)
		}
		return fmt.Errorf("Unable to mine while the ship's health is zero")
	}

	query = db.sqlx.Rebind(fmt.Sprintf(`
		UPDATE %s SET
			remaining = remaining - ?,
			updated_time = now()
		WHERE id = ?
	`, asteroidsTable))
	_, err = tx.Exec(query, shares*ResourceToShareRatio, asteroid.Id)
	if err != nil {
		return err
	}

	query = db.sqlx.Rebind(fmt.Sprintf(`
		UPDATE %s AS ships SET
			health = health - ?
		FROM %s AS ast
		WHERE ast.id = ? AND ast.ship_id = ships.id
	`, shipsTable, asteroidsTable))
	incr := int64(math.Min(float64(ship.Health), elapsedTime*float64(damagePerSec)))
	_, err = tx.Exec(query, incr, asteroid.Id)

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

func (db *store) AssignAsteroid(id int64, ship Ship) error {
	tx, err := db.sqlx.Beginx()
	if err != nil {
		return err
	}

	var size int
	query := db.sqlx.Rebind(
		fmt.Sprintf(`SELECT total FROM %s WHERE id = ?`, asteroidsTable),
	)
	err = tx.Get(&size, query, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if size > ship.Cargo {
		return fmt.Errorf("This asteroid is too large for your cargo hold.")
	}

	query = db.sqlx.Rebind(fmt.Sprintf(`
        UPDATE %s AS ast SET
            ship_id         = ?,
			ship_speed		= ?,
            updated_time    = now()
        WHERE ast.id = ? AND ast.ship_id = 0`,
		asteroidsTable))
	_, err = tx.Exec(query, ship.Id, ship.Speed, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *store) DestroyAsteroids() error {
	query := fmt.Sprintf("DELETE FROM %s WHERE remaining < 1", asteroidsTable)
	_, err := db.sqlx.Exec(db.sqlx.Rebind(query))
	return err
}
