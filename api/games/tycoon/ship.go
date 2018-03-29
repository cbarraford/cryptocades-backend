package tycoon

import (
	"errors"
	"fmt"
	"log"
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
		userId, err := context.GetUserId(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		shipId, err := context.GetInt64("id", c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var json inputShip
		err = c.BindJSON(&json)
		if err != nil {

			log.Printf("Error: %+v", err)
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
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
			return
		}
		if shipUserId != userId {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("Ship Access Denied"))
		}
		log.Print("FOO")
		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_ships",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		ship, err := store.GetShip(shipId)
		seg.End()
		if err != nil {
			log.Printf("FOOBAR %+v", err)
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
