package users

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cbarraford/cryptocades-backend/api/context"
	"github.com/cbarraford/cryptocades-backend/store/user"
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
		if err != nil {

			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		if json.BTCAddr != "" {
			record.BTCAddr = json.BTCAddr
			err = store.Update(&record)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}

		// check if password is being changed
		if json.Password != "" {
			record.Password = json.Password
			err = store.PasswordSet(&record)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}

		c.JSON(http.StatusOK, record)
	}
}
