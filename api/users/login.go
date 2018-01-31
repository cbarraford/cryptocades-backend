package users

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/store/session"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(store user.Store, sessionStore session.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		txn := nrgin.Transaction(c)
		record := user.Record{}

		var json login
		err = c.BindJSON(&json)
		if err == nil {
			seg := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: "users",
				Operation:  "Auth",
			}
			seg.StartTime = newrelic.StartSegmentNow(txn)
			record, err = store.Authenticate(json.Username, json.Password)
			seg.End()
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
		} else {
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		sessionRecord := session.Record{
			UserId: record.Id,
		}
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "session",
			Operation:  "Create",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = sessionStore.Create(&sessionRecord, 30)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, sessionRecord)
	}
}
