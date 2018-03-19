package admins

import (
	"net/http"
	"strconv"

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

func TotalActiveUsers(store admin.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		min := c.DefaultQuery("minutes", "1440")
		minutes, err := strconv.Atoi(min)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "users",
			Operation:  "ActiveUsers",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		active, err := store.TotalActiveUsers(minutes)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		seg.End()
		c.JSON(http.StatusOK, gin.H{"active": active})
	}
}

func TotalLiveUsers(store admin.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		min := c.DefaultQuery("minutes", "5")
		minutes, err := strconv.Atoi(min)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "users",
			Operation:  "LiveUsers",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		live, err := store.TotalActiveUsers(minutes)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		seg.End()
		c.JSON(http.StatusOK, gin.H{"live": live})
	}
}

type inputPlays struct {
	Email  string `json:"email"`
	Title  string `json:"title"`
	Amount int    `json:"amount"`
}

func AwardPlays(store admin.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		var json inputPlays
		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "incomes",
			Operation:  "AwardPlays",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = store.AwardPlays(json.Email, json.Amount, json.Title)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		seg.End()
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	}
}
