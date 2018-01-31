package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/api/context"
	"github.com/cbarraford/cryptocades-backend/store/confirmation"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

func Confirm(confirmStore confirmation.Store, store user.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		code := context.GetString("code", c)

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "confirmation",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		confirm, err := confirmStore.GetByCode(code)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "users",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		user, err := store.Get(confirm.UserId)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		// if the confirmation is email change, update email
		user.Email = confirm.Email

		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "users",
			Operation:  "Confirmed",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = store.MarkAsConfirmed(&user)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "confirmations",
			Operation:  "Delete",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = confirmStore.Delete(confirm.Id)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	}
}
