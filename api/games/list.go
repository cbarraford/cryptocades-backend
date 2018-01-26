package games

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cbarraford/cryptocades-backend/store/game"
)

func List(store game.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, store.List())
	}
}
