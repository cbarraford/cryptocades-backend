package jackpots

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CBarraford/lotto/store/jackpot"
)

func List(store jackpot.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		records, err := store.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, records)
		}
	}
}
