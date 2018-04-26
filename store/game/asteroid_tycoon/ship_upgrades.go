package asteroid_tycoon

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

type ShipUpgrade struct {
	CategoryId  int       `json:"category_id" db:"category_id"`
	AssetId     int       `json:"asset_id" db:"asset_id"`
	Value       int       `json:"value" db:"value"`
	Cost        int       `json:"cost" db:"cost"`
	Id          int64     `json:"id" db:"id"`
	ShipId      int64     `json:"ship_id" db:"ship_id"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
}

const upgradesTable string = "g2_applied_ship_upgrades"
const listUpgradesTable string = "g2_ship_upgrades"

func (db *store) InitShip(shipId int64) error {
	var err error

	query := db.sqlx.Rebind(fmt.Sprintf(`
		INSERT INTO %s
			(ship_id, category_id, asset_id)
		VALUES
			(?, ?, ?)
	`, upgradesTable))

	_, err = db.sqlx.Exec(query, shipId, 1, 1)
	if err != nil {
		return err
	}
	_, err = db.sqlx.Exec(query, shipId, 2, 1)
	if err != nil {
		return err
	}
	_, err = db.sqlx.Exec(query, shipId, 3, 1)
	if err != nil {
		return err
	}
	_, err = db.sqlx.Exec(query, shipId, 4, 1)
	if err != nil {
		return err
	}

	return nil

}

func (db *store) ApplyUpgrade(shipId int64, upgrade ShipUpgrade) error {
	var err error
	// check for required attrs
	if shipId == 0 {
		return fmt.Errorf("Must belong to a ship.")
	}
	if upgrade.CategoryId == 0 {
		return fmt.Errorf("Must be in a category.")
	}
	if upgrade.AssetId == 0 {
		return fmt.Errorf("Must have an asset id.")
	}
	tx, err := db.sqlx.Beginx()

	var assetId int
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT asset_id FROM %s WHERE ship_id = ? AND category_id = ?", upgradesTable))
	err = tx.Get(&assetId, query, shipId, upgrade.CategoryId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if assetId != upgrade.AssetId-1 {
		tx.Rollback()
		return fmt.Errorf("Upgrades can only be applied in sequential order")
	}

	// subtract the cost of the upgrade
	query = db.sqlx.Rebind(fmt.Sprintf(`
        UPDATE %s AS acct SET  
            credits			= credits - ?,
            updated_time	= now()
        FROM %s AS ship
		WHERE ship.account_id = acct.Id AND ship.Id = ?`, accountsTable, shipsTable))
	_, err = tx.Exec(query, upgrade.Cost, shipId)
	if serr, ok := err.(*pq.Error); ok {
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
		UPDATE %s
			set asset_id = ?
		WHERE ship_id = ? AND category_id = ? AND asset_id = ?`, upgradesTable))
	_, err = tx.Exec(
		query,
		upgrade.AssetId,
		shipId,
		upgrade.CategoryId,
		upgrade.AssetId-1,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *store) GetUpgrade(categoryId, assetId int) (ShipUpgrade, error) {
	var err error
	upgrade := ShipUpgrade{}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE category_id = ? AND asset_id = ?", listUpgradesTable))
	err = db.sqlx.Get(&upgrade, query, categoryId, assetId)
	return upgrade, err
}

func (db *store) ListUpgrades() ([]ShipUpgrade, error) {
	var err error
	upgrades := []ShipUpgrade{}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s", listUpgradesTable))
	err = db.sqlx.Select(&upgrades, query)

	return upgrades, err
}

func (db *store) GetUpgradesByShipId(shipId int64) ([]ShipUpgrade, error) {
	var err error
	upgrades := []ShipUpgrade{}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE ship_id = ?", upgradesTable))
	err = db.sqlx.Select(&upgrades, query, shipId)

	return upgrades, err
}
