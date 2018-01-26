package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cbarraford/cryptocades-backend/api/context"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

func Get(store user.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := context.GetInt64("id", c)
		if err != nil {
			return
		}

		record, err := store.Get(id)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// obscure email address
		record.Email = ""

		c.JSON(http.StatusOK, record)
	}
}
