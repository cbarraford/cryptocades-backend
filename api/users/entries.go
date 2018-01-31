package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/api/context"
	"github.com/cbarraford/cryptocades-backend/store/entry"
)

func Entries(entries entry.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		var userId int64

		if userId, err = context.GetUserId(c); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "entries",
			Operation:  "LIST",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		records, err := entries.ListByUser(userId)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, records)
	}
}
