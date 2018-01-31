package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/api/context"
	"github.com/cbarraford/cryptocades-backend/store/entry"
	"github.com/cbarraford/cryptocades-backend/store/income"
)

func Balance(store income.Store, entries entry.Store) func(*gin.Context) {
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
			Collection: "incomes",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		income, err := store.UserIncome(userId)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		txn = nrgin.Transaction(c)
		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "entries",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		spent, err := entries.UserSpent(userId)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"balance": income - spent})
	}
}
