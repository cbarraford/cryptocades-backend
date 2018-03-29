package asteroid_tycoon

import "errors"

var kaboom = errors.New("Not Implemented")

type Dummy struct{}

// Account
func (*Dummy) CreateAccount(acct *Account) error                { return kaboom }
func (*Dummy) GetAccountByUserId(userId int64) (Account, error) { return Account{}, kaboom }
func (*Dummy) UpdateAccount(acct *Account) error                { return kaboom }
func (*Dummy) DeleteAccount(id int64) error                     { return kaboom }

// Ship
func (*Dummy) CreateShipt(ship *Ship) error                     { return kaboom }
func (*Dummy) GetShipsByAccountId(acctId int64) ([]Ship, error) { return nil, kaboom }
func (*Dummy) GetShipUserId(id int64) (int64, error)            { return 0, kaboom }
func (*Dummy) GetShip(id int64) (Ship, error)                   { return Ship{}, kaboom }
func (*Dummy) UpdateShip(ship *Ship) error                      { return kaboom }
func (*Dummy) AddResources(a, r int) error                      { return kaboom }
func (*Dummy) AddDamage(h, d int) error                         { return kaboom }
func (*Dummy) DeleteShip(id int64) error                        { return kaboom }

// Upgrades
func (*Dummy) ApplyUpgrade(up *AppliedUpgrade) error { return kaboom }

// Asteroids
func (*Dummy) CreateAsteroid(ast *Asteroid) error              { return kaboom }
func (*Dummy) AssignAsteroid(id int64, ship *Ship) error       { return kaboom }
func (*Dummy) OwnedAsteroids(shipId int64) ([]Asteroid, error) { return nil, kaboom }
func (*Dummy) AvailableAsteroids() ([]Asteroid, error)         { return nil, kaboom }
func (*Dummy) DestroyAsteroids() error                         { return kaboom }

// Logs
func (*Dummy) CreateLog(line *Log) error               { return kaboom }
func (*Dummy) GetShipLogs(shipId int64) ([]Log, error) { return nil, kaboom }
func (*Dummy) DeleteOldLogs() error                    { return kaboom }