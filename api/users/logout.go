package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/store/session"
)

func Logout(store session.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		token := c.Request.Header["Session"][0]
		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "sessions",
			Operation:  "DELETE",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = store.Delete(token)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}
}
