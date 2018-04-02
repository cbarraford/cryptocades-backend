package tycoon

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/api/context"
	"github.com/cbarraford/cryptocades-backend/store/game/asteroid_tycoon"
)

func CreateAccount(store asteroid_tycoon.Store) func(*gin.Context) {
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
			Collection: "g2_accounts",
			Operation:  "CREATE",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		account := asteroid_tycoon.Account{
			UserId: userId,
		}
		err = store.CreateAccount(&account)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, account)
		}
	}
}

func GetAccount(store asteroid_tycoon.Store) func(*gin.Context) {
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
			Collection: "g2_accounts",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		account, err := store.GetAccountByUserId(userId)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, account)
		}
	}
}

type tradeInput struct {
	Amount int `json:"amount"`
}

func TradeForCredits(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		userId, err := context.GetUserId(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var json tradeInput
		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_accounts",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		account, err := store.GetAccountByUserId(userId)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_accounts",
			Operation:  "TRADE",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = store.TradeForCredits(account.Id, json.Amount)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "OK"})
		}
	}
}

func TradeForPlays(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		userId, err := context.GetUserId(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var json tradeInput
		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_accounts",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		account, err := store.GetAccountByUserId(userId)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "g2_accounts",
			Operation:  "TRADE",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = store.TradeForPlays(account.Id, json.Amount)
		seg.End()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "OK"})
		}
	}
}
