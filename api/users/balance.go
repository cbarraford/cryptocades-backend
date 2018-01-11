package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CBarraford/lotto/api/context"
	"github.com/CBarraford/lotto/store/entry"
	"github.com/CBarraford/lotto/store/user"
)

func Balance(store user.Store, entries entry.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		var userId int64

		if userId, err = context.GetUserId(c); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		record, err := store.Get(userId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		spent, err := entries.UserSpent(userId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		balance := (record.MinedHashes + record.BonusHashes) - spent

		c.JSON(http.StatusOK, gin.H{"balance": balance})
	}
}
