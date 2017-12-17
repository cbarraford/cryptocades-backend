package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CBarraford/lotto/store/session"
)

type logout struct {
	Token string `json:"token"`
}

func Logout(store session.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		var json logout
		err = c.BindJSON(&json)
		if err == nil {
			err = store.Delete(json.Token)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
		} else {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Could not parse json body: %+v", err))
			return
		}
	}
}
