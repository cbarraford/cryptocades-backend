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
	State          int       `json:"-" db:"state"`
	TotalAsteroids int       `json:"total_asteroids" db:"total_asteroids"`
	TotalResources int       `json:"total_resources" db:"total_resources"`
	Health         int       `json:"health" db:"health"`
	DrillBit       int       `json:"drill_bit" db:"drill_bit"`
	SolarSystem    int       `json:"-" db:"solar_system"`
	CreatedTime    time.Time `json:"created_time" db:"created_time"`
	UpdatedTime    time.Time `json:"updated_time" db:"updated_time"`
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
	return err
}

func (db *store) GetShipsByAccountId(accountId int64) ([]Ship, error) {
	var err error
	ships := []Ship{}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE account_id = ?", shipsTable))
	err = db.sqlx.Select(&ships, query, accountId)
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

func (db *store) AddResources(asteroids, resources int) error {
	query := db.sqlx.Rebind(fmt.Sprintf(
		"UPDATE %s SET total_asteroids = total_asteroids + ?, total_resources = total_resources + ?", shipsTable,
	))
	_, err := db.sqlx.Exec(query, asteroids, resources)
	return err
}

func (db *store) AddDamage(health, drillbit int) error {
	query := db.sqlx.Rebind(fmt.Sprintf(
		"UPDATE %s SET health = health - ?, drill_bit = drill_bit - ?", shipsTable,
	))
	_, err := db.sqlx.Exec(query, health, drillbit)
	return err
}

func (db *store) DeleteShip(id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", shipsTable)
	_, err := db.sqlx.Exec(db.sqlx.Rebind(query), id)
	return err
}
