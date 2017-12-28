package users

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CBarraford/lotto/api/context"
	"github.com/CBarraford/lotto/store/user"
)

func Update(store user.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		userId, err := context.GetUserId(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		record, err := store.Get(userId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var json input
		err = c.BindJSON(&json)
		if err == nil {
			record.Username = json.Username
			record.Email = json.Email
			record.BTCAddr = json.BTCAddr
		} else {
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		err = store.Update(&record)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, record)
	}
}
