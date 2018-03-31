package asteroid_tycoon

import (
	"fmt"
	"log"
	"time"
)

type Category struct {
	Id          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Upgrades    []ShipUpgrade `json:"upgrades"`
}

type ShipUpgrade struct {
	CategoryId int    `json:"category_id"`
	AssetId    int    `json:"asset_id"`
	Value      int    `json:"value"`
	Name       string `json:"name"`
	Cost       int    `json:"cost"`
	//Description string `json:"description"`
}

type AppliedUpgrade struct {
	Id          int64     `json:"id" db:"id"`
	ShipId      int64     `json:"-" db:"ship_id"`
	CategoryId  int       `json:"-" db:"category_id"`
	AssetId     int       `json:"-" db:"asset_id"`
	CreatedTime time.Time `json:"created_time" db:"created_time"`
}

var Categories = []Category{
	{
		1,
		"Engines",
		"Engines push your ship through space. Faster engines means faster top speed",
		[]ShipUpgrade{
			{1, 1, 100, "Basic Engine", 100},
		},
	},
	{
		2,
		"Cargo",
		"Cargo is where you store your mined resources. A larger cargo allows you to mine larger asteroids as well as make less trips back and forth to the hanger.",
		[]ShipUpgrade{
			{2, 1, 100, "500 Cargo", 100},
		},
	},
	{
		3,
		"Drill",
		"Drill is what collects resources from asteroids. The better the drill, the longer your drill bits last before they need to be replaced.",
		[]ShipUpgrade{
			{3, 1, 100, "Copper Drill", 100},
		},
	},
	{
		4,
		"Hull",
		"The hull protects your from the debris and raditation that damages your ship over time. A stronger hull ensures you can stay out longer collecting those valuable resources",
		[]ShipUpgrade{
			{4, 1, 100, "Copper Hull", 100},
			{4, 2, 200, "Aluminium Hull", 200},
			{4, 3, 300, "Iron Hull", 300},
			{4, 4, 400, "Steel Hull", 400},
			{4, 5, 500, "Titanium Hull", 500},
			{4, 6, 600, "Titanium Aluminide Hull", 600},
			{4, 7, 700, "Tungsten Hull", 700},
			{4, 8, 800, "Tungsten Carbide Hull", 800},
			{4, 9, 900, "Inconel Hull", 900},
			{4, 10, 1000, "Carbon Steel", 1000},
		},
	},
}

const upgradesTable string = "g2_ship_upgrades"

func (db *store) InitShip(shipId int64) error {

	query := db.sqlx.Rebind(fmt.Sprintf(`
		INSERT INTO %s
			(ship_id, category_id, asset_id)
		VALUES
			(?, ?, ?)
	`, upgradesTable))

	for _, cat := range Categories {
		up := cat.Upgrades[0]
		_, err := db.sqlx.Exec(query, shipId, up.CategoryId, up.AssetId)
		if err != nil {
			return err
		}
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
            credits		= credits - ?,
            updated_time   = now()
        FROM %s AS ship
		WHERE ship.account_id = acct.Id AND ship.Id = ?`, accountsTable, shipsTable))
	_, err = tx.Exec(query, upgrade.Cost, shipId)
	if err != nil {
		tx.Rollback()
		log.Printf("Not enough credits: %s", err)
		return fmt.Errorf("Not enough credits.")
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

func (db *store) GetUpgradesByShipId(shipId int64) ([]AppliedUpgrade, error) {
	var err error
	upgrades := []AppliedUpgrade{}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE ship_id = ?", upgradesTable))
	err = db.sqlx.Select(&upgrades, query, shipId)
	return upgrades, err
}
