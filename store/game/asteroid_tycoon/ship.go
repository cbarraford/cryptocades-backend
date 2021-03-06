package asteroid_tycoon

import (
	"fmt"
	"time"
)

const shipsTable string = "g2_ships"

type Ship struct {
	Id             int64     `json:"id" db:"id"`
	AccountId      int64     `json:"account_id" db:"account_id"`
	Name           string    `json:"name" db:"name"`
	TotalAsteroids int       `json:"total_asteroids" db:"total_asteroids"`
	TotalResources int       `json:"total_resources" db:"total_resources"`
	Health         int       `json:"health" db:"health"`
	SolarSystem    int       `json:"-" db:"solar_system"`
	Speed          int       `json:"speed" db:"-"`
	Hull           int       `json:"hull" db:"-"`
	Cargo          int       `json:"cargo" db:"-"`
	Repair         int       `json:"repair" db:"-"`
	SessionId      string    `json:"-" db:"session_id"`
	CreatedTime    time.Time `json:"created_time" db:"created_time"`
	UpdatedTime    time.Time `json:"updated_time" db:"updated_time"`
}

type ShipStatus struct {
	Status        string   `json:"status"`
	RemainingTime int      `json:"remaining_time"`
	TravelTime    int      `json:"travel_time"`
	Asteroid      Asteroid `json:"asteroid"`
	Ship          Ship     `json:"ship"`
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

	if ship.Health == 0 {
		ship.Health = 100
	}

	query := fmt.Sprintf(`
		INSERT INTO %s
			(account_id, name, total_asteroids, total_resources, health, created_time)
		VALUES
			(:account_id, :name, :total_asteroids, :total_resources, :health, :created_time) RETURNING id
	`, shipsTable)

	stmt, err := db.sqlx.PrepareNamed(query)
	err = stmt.QueryRowx(ship).Scan(&ship.Id)
	if err != nil {
		return err
	}

	err = db.InitShip(ship.Id)
	if err != nil {
		return err
	}

	return db.ExpandShip(ship)
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
			total_asteroids = :total_asteroids,
			total_resources = :total_resources,
			health			= :health,
			session_id		= :session_id,
			solar_system	= :solar_system,
            updated_time    = :updated_time
        WHERE id = :id`, shipsTable)
	_, err := db.sqlx.NamedExec(query, ship)
	return err
}

func (db *store) ExpandShip(ship *Ship) error {
	// TODO: Build cache here, since the data here doesn't change after boot,
	// no need to keep querying the db and add additional load
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
	err = db.sqlx.Get(&ship.Repair, query, ship.Id, 3)
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

func (db *store) GetStatus(ship Ship, ast Asteroid) (status ShipStatus) {
	status.Asteroid = ast
	status.Ship = ship
	status.Status = "Docked"
	if ast.Id == 0 {
		return
	}

	status.TravelTime = status.Asteroid.Distance
	if status.Asteroid.ShipSpeed > 0 {
		status.TravelTime = status.Asteroid.Distance / status.Asteroid.ShipSpeed
	}
	diffTime := time.Now().Unix() - status.Asteroid.UpdatedTime.Unix()
	status.RemainingTime = status.TravelTime - int(diffTime)

	if ship.Id == 0 || ship.Health <= 0 {
		return
	}

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
