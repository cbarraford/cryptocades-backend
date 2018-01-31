package jackpots

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/api/context"
	"github.com/cbarraford/cryptocades-backend/store/entry"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

type input struct {
	Amount int `json:"amount"`
}

func Enter(store entry.Store, userStore user.Store) func(*gin.Context) {
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
			Collection: "users",
			Operation:  "LIST",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		user, err := userStore.Get(userId)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		btcRegex, _ := regexp.Compile("^[13][a-km-zA-HJ-NP-Z0-9]{26,33}$")
		if !btcRegex.MatchString(user.BTCAddr) {
			err = fmt.Errorf("Must have a valid bitcoin address to enter a jackpot.")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		jackpotId, err := context.GetInt64("id", c)
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

		record := entry.Record{
			JackpotId: jackpotId,
			UserId:    userId,
			Amount:    json.Amount,
		}

		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "entries",
			Operation:  "Create",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = store.Create(&record)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	}
}
