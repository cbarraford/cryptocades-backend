package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CBarraford/lotto/store/session"
)

func Logout(store session.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		token := c.Request.Header["Session"][0]
		err = store.Delete(token)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}
}
