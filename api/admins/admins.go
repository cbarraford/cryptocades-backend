package admins

import (
	"net/http"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	admin "github.com/cbarraford/cryptocades-backend/admin"
)

func TotalRegisterUsers(store admin.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "users",
			Operation:  "TotalUsers",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		total, err := store.TotalRegisterUsers()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		seg.End()
		c.JSON(http.StatusOK, gin.H{"total": total})
	}
}
