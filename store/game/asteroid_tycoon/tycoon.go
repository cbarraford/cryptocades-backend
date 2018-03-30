package asteroid_tycoon

import "github.com/jmoiron/sqlx"

type Store interface {
	// Account
	CreateAccount(acct *Account) error
	GetAccountByUserId(userId int64) (Account, error)
	UpdateAccount(acct *Account) error
	DeleteAccount(id int64) error

	// Ship
	CreateShip(ship *Ship) error
	GetShipsByAccountId(acctId int64) ([]Ship, error)
	GetShipUserId(id int64) (int64, error)
	GetShip(id int64) (Ship, error)
	UpdateShip(ship *Ship) error
	AddResources(a, r int) error
	AddDamage(h, d int) error
	DeleteShip(id int64) error

	// Upgrades
	ApplyUpgrade(up *AppliedUpgrade) error
	GetUpgradesByShipId(shipId int64) ([]AppliedUpgrade, error)

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
