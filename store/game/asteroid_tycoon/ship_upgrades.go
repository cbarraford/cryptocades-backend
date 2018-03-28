package asteroid_tycoon

import (
	"fmt"
	"time"
)

type Category struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ShipUpgrade struct {
	CategoryId int    `json:"category_id"`
	AssetId    int    `json:"asset_id"`
	Value      int    `json:"value"`
	Name       string `json:"name"`
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
	},
	{
		2,
		"Cargo",
		"Cargo is where you store your mined resources. A larger cargo allows you to mine larger asteroids as well as make less trips back and forth to the hanger.",
	},
	{
		3,
		"Drill",
		"Drill is what collects resources from asteroids. The better the drill, the longer your drill bits last before they need to be replaced.",
	},
	{
		4,
		"Hull",
		"The hull protects your from the debris and raditation that damages your ship over time. A stronger hull ensures you can stay out longer collecting those valuable resources",
	},
}

var ShipUpgrades = []ShipUpgrade{
	{1, 1, 100, "Copper Hull"},
	{1, 2, 200, "Aluminium Hull"},
	{1, 3, 200, "Iron Hull"},
	{1, 4, 400, "Steel Hull"},
	{1, 5, 500, "Titanium Hull"},
	{1, 6, 600, "Titanium Aluminide Hull"},
	{1, 7, 700, "Tungsten Hull"},
	{1, 8, 800, "Tungsten Carbide Hull"},
	{1, 9, 900, "Inconel Hull"},
	{1, 10, 1000, "Carbon Steel"},
}

const upgradesTable string = "g2_ship_upgrades"

func (db *store) ApplyUpgrade(upgrade *AppliedUpgrade) error {
	var err error
	// check for required attrs
	if upgrade.ShipId == 0 {
		return fmt.Errorf("Must belong to a ship.")
	}
	if upgrade.CategoryId == 0 {
		return fmt.Errorf("Must be in a category.")
	}
	if upgrade.AssetId == 0 {
		return fmt.Errorf("Must have an asset id.")
	}

	query := fmt.Sprintf(`
		INSERT INTO %s 
			(ship_id, category_id, asset_id) 
		VALUES 
			(?,?,?) ON CONFLICT (ship_id, category_id) DO UPDATE SET asset_id = ?;`, upgradesTable)
	query = db.sqlx.Rebind(query)
	_, err = db.sqlx.Exec(
		query,
		upgrade.ShipId,
		upgrade.CategoryId,
		upgrade.AssetId,
		upgrade.AssetId,
	)
	return err
}

func (db *store) GetUpgradesByShipId(shipId int64) ([]AppliedUpgrade, error) {
	var err error
	upgrades := []AppliedUpgrade{}
	query := db.sqlx.Rebind(fmt.Sprintf("SELECT * FROM %s WHERE ship_id = ?", upgradesTable))
	err = db.sqlx.Select(&upgrades, query, shipId)
	return upgrades, err
}
