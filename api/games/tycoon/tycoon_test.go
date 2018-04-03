package tycoon

import (
	"testing"

	check "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store/game/asteroid_tycoon"
)

func TestPackage(t *testing.T) { check.TestingT(t) }

type mockStore struct {
	asteroid_tycoon.Dummy
	created    bool
	updated    bool
	name       string
	userId     int64
	accountId  int64
	categoryId int
	assetId    int
	asteroidId int64
	shipId     int64
	amount     int
}

func (m *mockStore) GetAccountByUserId(userId int64) (asteroid_tycoon.Account, error) {
	return asteroid_tycoon.Account{
		Id:     1,
		UserId: userId,
	}, nil
}

func (m *mockStore) CreateAccount(acct *asteroid_tycoon.Account) error {
	m.created = true
	m.userId = acct.UserId
	return nil
}

func (m *mockStore) CreateShip(ship *asteroid_tycoon.Ship) error {
	m.created = true
	m.accountId = ship.AccountId
	return nil
}

func (m *mockStore) GetShipUserId(id int64) (int64, error) {
	return m.userId, nil
}

func (m *mockStore) GetShipsByAccountId(id int64) ([]asteroid_tycoon.Ship, error) {
	return []asteroid_tycoon.Ship{
		{AccountId: id},
	}, nil
}

func (m *mockStore) GetShip(id int64) (asteroid_tycoon.Ship, error) {
	return asteroid_tycoon.Ship{
		Id:        id,
		AccountId: m.accountId,
	}, nil
}

func (m *mockStore) UpdateShip(ship *asteroid_tycoon.Ship) error {
	m.updated = true
	m.name = ship.Name
	return nil
}

func (m *mockStore) GetShipLogs(shipId int64) ([]asteroid_tycoon.Log, error) {
	return []asteroid_tycoon.Log{
		{Log: "log-line-text"},
	}, nil
}

func (m *mockStore) ApplyUpgrade(shipId int64, upgrade asteroid_tycoon.ShipUpgrade) error {
	m.updated = true
	m.categoryId = upgrade.CategoryId
	m.assetId = upgrade.AssetId
	return nil
}

func (m *mockStore) AssignAsteroid(id, shipId int64) error {
	m.created = true
	m.shipId = shipId
	m.asteroidId = id
	return nil
}

func (m *mockStore) OwnedAsteroid(shipId int64) (asteroid_tycoon.Asteroid, error) {
	return asteroid_tycoon.Asteroid{
		ShipId: shipId,
	}, nil
}

func (m *mockStore) AvailableAsteroids() ([]asteroid_tycoon.Asteroid, error) {
	return []asteroid_tycoon.Asteroid{
		{ShipId: 0},
	}, nil
}

func (m *mockStore) TradeForPlays(accountId int64, amount int) error {
	m.accountId = accountId
	m.amount = amount
	return nil
}

func (m *mockStore) TradeForCredits(accountId int64, amount int) error {
	m.accountId = accountId
	m.amount = amount
	return nil
}
