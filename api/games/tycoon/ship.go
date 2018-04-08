package tycoon

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/api/context"
	"github.com/cbarraford/cryptocades-backend/store/game/asteroid_tycoon"
)

type inputShip struct {
	Name string `json:"name" db:"name"`
}

type inputUpgrade struct {
	CategoryId int `json:"category_id" db:"category_id"`
	AssetId    int `json:"asset_id" db:"asset_id"`
}

func authShip(c *gin.Context, store asteroid_tycoon.Store) error {
	var err error
	userId, err := context.GetUserId(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return err
	}

	shipId, err := context.GetInt64("id", c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return err
	}

	txn := nrgin.Transaction(c)
	seg := newrelic.DatastoreSegment{
		Product:    newrelic.DatastorePostgres,
		Collection: "g2_ships",
		Operation:  "AUTH",
	}
	seg.StartTime = newrelic.StartSegmentNow(txn)
	shipUserId, err := store.GetShipUserId(shipId)
	seg.End()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}
	if shipUserId != userId {
		err = fmt.Errorf("Ship Access Denied")
		c.AbortWithError(http.StatusForbidden, err)
		return err
	}

	return nil
}

func CreateShip(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		userId, err := context.GetUserId(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_accounts",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		account, err := store.GetAccountByUserId(userId)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_ships",
			Operation:  "CREATE",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		ship := asteroid_tycoon.Ship{
			AccountId: account.Id,
		}
		err = store.CreateShip(&ship)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, ship)
		}
	}
}

func GetShips(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		userId, err := context.GetUserId(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_accounts",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		account, err := store.GetAccountByUserId(userId)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_ships",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		ships, err := store.GetShipsByAccountId(account.Id)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, ships)
		}
	}
}

func UpdateShip(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		shipId, err := context.GetInt64("id", c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var json inputShip
		err = c.BindJSON(&json)
		if err != nil {

			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		if err := authShip(c, store); err != nil {
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_ships",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		ship, err := store.GetShip(shipId)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		ship.Name = json.Name

		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_ships",
			Operation:  "UPDATE",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = store.UpdateShip(&ship)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, ship)
		}
	}
}

func GetShipLogs(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		if err := authShip(c, store); err != nil {
			return
		}

		shipId, err := context.GetInt64("id", c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_logs",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		lines, err := store.GetShipLogs(shipId)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, lines)
		}
	}
}

func GetShipUpgrades(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		if err := authShip(c, store); err != nil {
			return
		}

		shipId, err := context.GetInt64("id", c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_logs",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		lines, err := store.GetUpgradesByShipId(shipId)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, lines)
		}
	}
}

func ApplyUpgrade(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		shipId, err := context.GetInt64("id", c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var json inputUpgrade
		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		if err := authShip(c, store); err != nil {
			return
		}

		up, err := store.GetUpgrade(json.CategoryId, json.AssetId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		up.ShipId = shipId

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_applied_ship_upgrades",
			Operation:  "UPDATE",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = store.ApplyUpgrade(shipId, up)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		} else {
			c.JSON(http.StatusOK, up)
			return
		}
	}
}

func GetMyAsteroids(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		shipId, err := context.GetInt64("id", c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if err := authShip(c, store); err != nil {
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_asteroids",
			Operation:  "LIST",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		asteroid, err := store.OwnedAsteroid(shipId)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, asteroid)
		}
	}
}

func Heal(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		shipId, err := context.GetInt64("id", c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if err := authShip(c, store); err != nil {
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_ships",
			Operation:  "HEAL",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = store.Heal(shipId)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "OK"})
		}
	}
}

func ReplaceDrillBit(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		shipId, err := context.GetInt64("id", c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if err := authShip(c, store); err != nil {
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_ships",
			Operation:  "REPLACE DRILLBIT",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = store.ReplaceDrillBit(shipId)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "OK"})
		}
	}
}

func GetStatus(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		shipId, err := context.GetInt64("id", c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if err := authShip(c, store); err != nil {
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_asteroids",
			Operation:  "LIST",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		asteroid, err := store.OwnedAsteroid(shipId)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, store.GetStatus(asteroid))
	}
}
