package jackpots

import (
	"net/http"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/store/jackpot"
)

func List(store jackpot.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "jackpots",
			Operation:  "LIST",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		records, err := store.List()
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, records)
		}
	}
}
