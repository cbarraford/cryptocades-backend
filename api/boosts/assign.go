package boosts

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/store/boost"
)

type input struct {
	BoostId  int64 `json:"boost_id"`
	IncomeId int64 `json:"income_id"`
}

func Assign(store boost.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		var json input
		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "boosts",
			Operation:  "Assign",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = store.Assign(json.BoostId, json.IncomeId)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	}
}
