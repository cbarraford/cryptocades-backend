package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CBarraford/lotto/api/context"
	"github.com/CBarraford/lotto/store/confirmation"
	"github.com/CBarraford/lotto/store/user"
)

func Confirm(confirmStore confirmation.Store, store user.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		code := context.GetString("code", c)

		confirm, err := confirmStore.GetByCode(code)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		user, err := store.Get(confirm.UserId)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		// if the confirmation is email change, update email
		user.Email = confirm.Email

		err = store.MarkAsConfirmed(&user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		err = confirmStore.Delete(confirm.Id)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	}
}
