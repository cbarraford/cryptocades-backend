package jackpots

import (
	"net/http"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/api/context"
	"github.com/cbarraford/cryptocades-backend/store/entry"
)

func Odds(store entry.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		var userId int64

		if userId, err = context.GetUserId(c); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		id, err := context.GetInt64("id", c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "entries",
			Operation:  "Odds",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		odds, err := store.GetOdds(id, userId)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, odds)
		}
	}
}
