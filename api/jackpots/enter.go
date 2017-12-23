package jackpots

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CBarraford/lotto/api/context"
	"github.com/CBarraford/lotto/store/entry"
	"github.com/CBarraford/lotto/store/user"
)

type input struct {
	Amount int `json:"amount"`
}

func Enter(userStore user.Store, store entry.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		var userId int64

		if userId, err = context.GetUserId(c); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
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

		// TODO: we should check total vs spent atomically

		user, err := userStore.Get(userId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		spent, err := store.UserSpent(userId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// ensure we're not spending currency we don't have
		if (spent + json.Amount) > (user.MinedHashes + user.BonusHashes) {
			err := fmt.Errorf("Insufficient funds.")
			c.AbortWithError(http.StatusPaymentRequired, err)
			return
		}

		record := entry.Record{
			JackpotId: jackpotId,
			UserId:    userId,
			Amount:    json.Amount,
		}

		err = store.Create(&record)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	}
}
