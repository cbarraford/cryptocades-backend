package asteroid_tycoon

import "errors"

var kaboom = errors.New("Not Implemented")

type Dummy struct{}

// Account
func (*Dummy) CreateAccount(acct *Account) error                   { return kaboom }
func (*Dummy) GetAccountByUserId(userId int64) (Account, error)    { return Account{}, kaboom }
func (*Dummy) UpdateAccount(acct *Account) error                   { return kaboom }
func (*Dummy) DeleteAccount(id int64) error                        { return kaboom }
func (*Dummy) AddAccountResources(id int64, amount int) error      { return kaboom }
func (*Dummy) SubtractAccountResources(id int64, amount int) error { return kaboom }
func (*Dummy) AddAccountCredits(id int64, amount int) error        { return kaboom }
func (*Dummy) SubtractAccountCredits(id int64, amount int) error   { return kaboom }
func (*Dummy) TradeForCredits(accountId int64, amount int) error   { return kaboom }
func (*Dummy) TradeForPlays(accountId int64, amount int) error     { return kaboom }

// Ship
func (*Dummy) CreateShipt(ship *Ship) error                     { return kaboom }
func (*Dummy) InitShip(shipId int64) error                      { return kaboom }
func (*Dummy) GetShipsByAccountId(acctId int64) ([]Ship, error) { return nil, kaboom }
func (*Dummy) GetShipUserId(id int64) (int64, error)            { return 0, kaboom }
func (*Dummy) GetShip(id int64) (Ship, error)                   { return Ship{}, kaboom }
func (*Dummy) UpdateShip(ship *Ship) error                      { return kaboom }
func (*Dummy) AddShipResources(a, r int) error                  { return kaboom }
func (*Dummy) AddShipDamage(h, d int) error                     { return kaboom }
func (*Dummy) DeleteShip(id int64) error                        { return kaboom }

// Upgrades
func (*Dummy) ApplyUpgrade(shipId int64, up ShipUpgrade) error {
	return kaboom
}
func (*Dummy) GetUpgradesByShipId(shipId int64) ([]AppliedUpgrade, error) { return nil, kaboom }

// Asteroids
func (*Dummy) CreateAsteroid(ast *Asteroid) error              { return kaboom }
func (*Dummy) AssignAsteroid(id, shipId int64) error           { return kaboom }
func (*Dummy) OwnedAsteroids(shipId int64) ([]Asteroid, error) { return nil, kaboom }
func (*Dummy) AvailableAsteroids() ([]Asteroid, error)         { return nil, kaboom }
func (*Dummy) DestroyAsteroids() error                         { return kaboom }

// Logs
func (*Dummy) CreateLog(line *Log) error               { return kaboom }
func (*Dummy) GetShipLogs(shipId int64) ([]Log, error) { return nil, kaboom }
func (*Dummy) DeleteOldLogs() error                    { return kaboom }
