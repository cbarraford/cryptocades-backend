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

	// Ship
	CreateShip(ship *Ship) error
	GetShipsByAccountId(acctId int64) ([]Ship, error)
	GetShipUserId(id int64) (int64, error)
	GetShip(id int64) (Ship, error)
	UpdateShip(ship *Ship) error
	AddShipResources(a, r int) error
	AddShipDamage(h, d int) error
	DeleteShip(id int64) error

	// Upgrades
	ApplyUpgrade(shipId int64, up ShipUpgrade) error
	GetUpgradesByShipId(shipId int64) ([]AppliedUpgrade, error)
	InitShip(shipId int64) error

	// Asteroids
	CreateAsteroid(ast *Asteroid) error
	AssignAsteroid(id, shipId int64) error
	OwnedAsteroids(shipId int64) ([]Asteroid, error)
	AvailableAsteroids() ([]Asteroid, error)
	DestroyAsteroids() error

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
