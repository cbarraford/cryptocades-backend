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

func CompletedAsteroid(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		var json inputAsteroid
		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		userId, err := context.GetUserId(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_ships",
			Operation:  "AUTH",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		shipUserId, err := store.GetShipUserId(json.ShipId)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		if shipUserId != userId {
			err = fmt.Errorf("Ship Access Denied")
			c.AbortWithError(http.StatusForbidden, err)
			return
		}

		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_asteroids",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		ast, err := store.OwnedAsteroid(json.ShipId)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_ledgers",
			Operation:  "CREATE",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = store.CompletedAsteroid(ast)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, gin.H{"Status": "OK"})
		}
	}
}
