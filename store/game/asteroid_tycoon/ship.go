package asteroid_tycoon

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

const shipsTable string = "g2_ships"

type Ship struct {
	Id             int64     `json:"id" db:"id"`
	AccountId      int64     `json:"account_id" db:"account_id"`
	Name           string    `json:"name" db:"name"`
	State          int       `json:"-" db:"state"`
	TotalAsteroids int       `json:"total_asteroids" db:"total_asteroids"`
	TotalResources int       `json:"total_resources" db:"total_resources"`
	Health         int       `json:"health" db:"health"`
	DrillBit       int       `json:"drill_bit" db:"drill_bit"`
	SolarSystem    int       `json:"-" db:"solar_system"`
	Speed          int       `json:"speed" db:"speed"`
	Hull           int       `json:"hull" db:"hull"`
	Cargo          int       `json:"cargo" db:"cargo"`
	Drill          int       `json:"drill" db:"drill"`
	CreatedTime    time.Time `json:"created_time" db:"created_time"`
	UpdatedTime    time.Time `json:"updated_time" db:"updated_time"`
}

type ShipStatus struct {
	Status        string   `json:"status"`
	RemainingTime int      `json:"remaining_time"`
	Asteroid      Asteroid `json:"asteroid"`
}

func (db *store) CreateShip(ship *Ship) error {
	var err error
	// check for required attrs
	if ship.AccountId == 0 {
		return fmt.Errorf("Must belong to an account.")
	}

	if ship.CreatedTime.IsZero() {
		ship.CreatedTime = time.Now().UTC()
	}

	if ship.Name == "" {
		ship.Name = "Eros 433"
	}

	query := fmt.Sprintf(`
		INSERT INTO %s
			(account_id, name, total_asteroids, total_resources, health, drill_bit, created_time)
		VALUES
			(:account_id, :name, :total_asteroids, :total_resources, :health, :drill_bit, :created_time) RETURNING id
	`, shipsTable)

	stmt, err := db.sqlx.PrepareNamed(query)
	err = stmt.QueryRowx(ship).Scan(&ship.Id)
	if err != nil {
		return err
	}

	return db.InitShip(ship.Id)
}

func (db *store) GetShipsByAccountId(accountId int64) ([]Ship, error) {
	var err error
	ships := []Ship{}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE account_id = ?", shipsTable))
	err = db.sqlx.Select(&ships, query, accountId)
	if err != nil {
		return ships, err
	}
	for i, _ := range ships {
		err = db.ExpandShip(&ships[i])
		if err != nil {
			return ships, err
		}
	}
	return ships, err
}

func (db *store) GetShipUserId(id int64) (int64, error) {
	var err error
	var userId int64
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT %s.user_id FROM %s JOIN %s ON %s.id = %s.account_id WHERE %s.id = ?", accountsTable, accountsTable, shipsTable, accountsTable, shipsTable, shipsTable))
	err = db.sqlx.Get(&userId, query, id)
	return userId, err
}

func (db *store) GetShip(id int64) (Ship, error) {
	var err error
	var ship Ship
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE id = ?", shipsTable))
	err = db.sqlx.Get(&ship, query, id)
	if err != nil {
		return ship, err
	}
	err = db.ExpandShip(&ship)
	return ship, err
}

func (db *store) UpdateShip(ship *Ship) error {
	// touch updated time
	ship.UpdatedTime = time.Now()

	query := fmt.Sprintf(`
        UPDATE %s SET
            name			= :name,
            state			= :state,
			total_asteroids = :total_asteroids,
			total_resources = :total_resources,
			health			= :health,
			drill_bit		= :drill_bit,
			solar_system	= :solar_system,
            updated_time    = :updated_time
        WHERE id = :id`, shipsTable)
	_, err := db.sqlx.NamedExec(query, ship)
	return err
}

func (db *store) Heal(shipId int64) error {
	ship, err := db.GetShip(shipId)
	if err != nil {
		return err
	}

	creditNeeded := (ship.Hull - ship.Health) / HealthForCredits

	tx, err := db.sqlx.Beginx()

	query := db.sqlx.Rebind(fmt.Sprintf(`
        UPDATE %s SET
			credits			= credits - ?,
            updated_time    = now()
        WHERE id = ?`, accountsTable))
	_, err = tx.Exec(query, creditNeeded, ship.AccountId)
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

	query = db.sqlx.Rebind(fmt.Sprintf(`
        UPDATE %s SET
			health			= ?,
            updated_time    = now()
        WHERE id = ?`, shipsTable))
	_, err = tx.Exec(query, ship.Hull, ship.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *store) ReplaceDrillBit(shipId int64) error {
	ship, err := db.GetShip(shipId)
	if err != nil {
		return err
	}

	var drillCost int
	query := db.sqlx.Rebind(fmt.Sprintf(`
		SELECT list.cost FROM %s AS ups JOIN %s AS list ON list.category_id = ups.category_id AND list.asset_id = ups.asset_id AND ups.ship_id = ? AND ups.category_id = ?`, upgradesTable, listUpgradesTable))
	err = db.sqlx.Get(&drillCost, query, ship.Id, 3)
	if err != nil {
		return err
	}
	creditNeeded := drillCost / DrillBitCost

	tx, err := db.sqlx.Beginx()

	query = db.sqlx.Rebind(fmt.Sprintf(`
        UPDATE %s SET
			credits			= credits - ?,
            updated_time    = now()
        WHERE id = ?`, accountsTable))
	_, err = tx.Exec(query, creditNeeded, ship.AccountId)
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

	query = db.sqlx.Rebind(fmt.Sprintf(`
        UPDATE %s SET
			drill_bit		= ?,
            updated_time    = now()
        WHERE id = ?`, shipsTable))
	_, err = tx.Exec(query, ship.Drill, ship.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *store) AddShipResources(asteroids, resources int) error {
	query := db.sqlx.Rebind(fmt.Sprintf(
		"UPDATE %s SET total_asteroids = total_asteroids + ?, total_resources = total_resources + ?", shipsTable,
	))
	_, err := db.sqlx.Exec(query, asteroids, resources)
	return err
}

func (db *store) AddShipDamage(health, drillbit int) error {
	query := db.sqlx.Rebind(fmt.Sprintf(
		"UPDATE %s SET health = health - ?, drill_bit = drill_bit - ?", shipsTable,
	))
	_, err := db.sqlx.Exec(query, health, drillbit)
	return err
}

func (db *store) ExpandShip(ship *Ship) error {
	var err error
	query := db.sqlx.Rebind(fmt.Sprintf(`
		SELECT list.value FROM %s AS ups JOIN %s AS list ON list.category_id = ups.category_id AND list.asset_id = ups.asset_id AND ups.ship_id = ? AND ups.category_id = ?`, upgradesTable, listUpgradesTable))
	err = db.sqlx.Get(&ship.Speed, query, ship.Id, 1)
	if err != nil {
		return err
	}
	err = db.sqlx.Get(&ship.Cargo, query, ship.Id, 2)
	if err != nil {
		return err
	}
	err = db.sqlx.Get(&ship.Drill, query, ship.Id, 3)
	if err != nil {
		return err
	}
	err = db.sqlx.Get(&ship.Hull, query, ship.Id, 4)
	if err != nil {
		return err
	}
	return nil
}

func (db *store) DeleteShip(id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", shipsTable)
	_, err := db.sqlx.Exec(db.sqlx.Rebind(query), id)
	return err
}

func (db *store) GetStatus(ast Asteroid) (status ShipStatus) {
	status.Asteroid = ast
	status.Status = "Docked"
	if ast.Id == 0 {
		return
	}

	travelTime := status.Asteroid.Distance
	if status.Asteroid.ShipSpeed > 0 {
		travelTime = status.Asteroid.Distance / status.Asteroid.ShipSpeed
	}
	diffTime := time.Now().Unix() - status.Asteroid.UpdatedTime.Unix()
	status.RemainingTime = travelTime - int(diffTime)
	if status.Asteroid.Remaining > 0 && status.Asteroid.Remaining < status.Asteroid.Total {
		status.Status = "Mining"
		return
	}
	if status.RemainingTime > 0 {
		if status.Asteroid.Remaining == status.Asteroid.Total {
			status.Status = "Approaching Asteroid"
		} else {
			status.Status = "Approaching Space Station"
		}
	} else {
		if status.Asteroid.Remaining == status.Asteroid.Total {
			status.Status = "Mining"
		} else {
			status.Status = "Docked"
		}
	}
	return
}
