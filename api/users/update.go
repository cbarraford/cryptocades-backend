package users

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/api/context"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

func Update(store user.Store) func(*gin.Context) {
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
			Collection: "users",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		record, err := store.Get(userId)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var json input
		err = c.BindJSON(&json)
		if err != nil {

			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		if json.BTCAddr != "" {
			record.BTCAddr = json.BTCAddr
			seg = newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: "users",
				Operation:  "Update",
			}
			seg.StartTime = newrelic.StartSegmentNow(txn)
			err = store.Update(&record)
			seg.End()
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}

		// check if password is being changed
		if json.Password != "" {
			record.Password = json.Password
			seg = newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: "users",
				Operation:  "Update",
			}
			seg.StartTime = newrelic.StartSegmentNow(txn)
			err = store.PasswordSet(&record)
			seg.End()
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}

		c.JSON(http.StatusOK, record)
	}
}
