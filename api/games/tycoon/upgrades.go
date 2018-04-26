package tycoon

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cbarraford/cryptocades-backend/store/game/asteroid_tycoon"
)

func GetUpgrades(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		upgrades, err := store.ListUpgrades()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, upgrades)
		}
	}
}
