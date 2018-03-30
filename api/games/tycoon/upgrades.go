package tycoon

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cbarraford/cryptocades-backend/store/game/asteroid_tycoon"
)

func GetUpgrades(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, asteroid_tycoon.ShipUpgrades)
	}
}
func GetCategories(store asteroid_tycoon.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, asteroid_tycoon.Categories)
	}
}
