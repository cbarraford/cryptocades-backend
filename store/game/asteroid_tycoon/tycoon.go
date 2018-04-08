package asteroid_tycoon

import "github.com/jmoiron/sqlx"

type Store interface {
	// Account
	CreateAccount(acct *Account) error
	GetAccountByUserId(userId int64) (Account, error)
	UpdateAccount(acct *Account) error
	AddAccountResources(id int64, amount int) error
	SubtractAccountResources(id int64, amount int) error
	AddAccountCredits(id int64, amount int) error
	SubtractAccountCredits(id int64, amount int) error
	DeleteAccount(id int64) error
	TradeForCredits(accountId int64, amount int) error
	TradeForPlays(accountId int64, amount int) error

	// Ship
	CreateShip(ship *Ship) error
	GetShipsByAccountId(acctId int64) ([]Ship, error)
	GetShipUserId(id int64) (int64, error)
	GetShip(id int64) (Ship, error)
	UpdateShip(ship *Ship) error
	Heal(shipId int64) error
	ReplaceDrillBit(shipId int64) error
	AddShipResources(a, r int) error
	AddShipDamage(h, d int) error
	DeleteShip(id int64) error
	GetStatus(ast Asteroid) (status ShipStatus)

	// Upgrades
	ApplyUpgrade(shipId int64, up ShipUpgrade) error
	GetUpgradesByShipId(shipId int64) ([]ShipUpgrade, error)
	InitShip(shipId int64) error
	ListUpgrades() ([]ShipUpgrade, error)
	GetUpgrade(categoryId, assetId int) (ShipUpgrade, error)

	// Asteroids
	CreateAsteroid(ast *Asteroid) error
	Mined(sessionId string, shares int, userId int64, tx *sqlx.Tx) error
	AssignAsteroid(id int64, sessionId string, ship Ship) error
	OwnedAsteroid(shipId int64) (Asteroid, error)
	AvailableAsteroids() ([]Asteroid, error)
	DestroyAsteroids() error

	// Ledger
	CompletedAsteroid(ast Asteroid) error
	AddEntry(accountId int64, amount int, description string) error
	ListByAccountId(accountId int64) ([]Ledger, error)

	// Logs
	CreateLog(line *Log) error
	GetShipLogs(shipId int64) ([]Log, error)
	DeleteOldLogs() error
}

type store struct {
	Store
	sqlx *sqlx.DB
}

func NewStore(db *sqlx.DB) Store {
	return &store{sqlx: db}
}
