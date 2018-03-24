package matchups

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/api/context"
	"github.com/cbarraford/cryptocades-backend/store/matchup"
)

func List(store matchup.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		top, err := strconv.Atoi(c.DefaultQuery("top", "20"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		event := context.GetString("event", c)

		offset, err := context.GetInt("offset", c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		offset = -offset

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "matchups",
			Operation:  "LIST",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		records, err := store.GetTopPerformers(event, offset, top)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		seg.End()
		c.JSON(http.StatusOK, records)
	}
}

func Get(store matchup.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		var userId int64

		if userId, err = context.GetUserId(c); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		event := context.GetString("event", c)

		offset, err := context.GetInt("offset", c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		offset = -offset

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "matchups",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		record, err := store.Get(event, offset, userId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		seg.End()
		c.JSON(http.StatusOK, record)
	}
}
