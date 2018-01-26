package users

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cbarraford/cryptocades-backend/api/context"
	"github.com/cbarraford/cryptocades-backend/store/confirmation"
	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/cbarraford/cryptocades-backend/util"
	emailer "github.com/cbarraford/cryptocades-backend/util/email"
	"github.com/cbarraford/cryptocades-backend/util/url"
)

type passwordEmail struct {
	Email string `json:"email"`
}
type passwordReset struct {
	Password string `json:"password"`
}

func PasswordReset(confirmStore confirmation.Store, store user.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		code := context.GetString("code", c)

		confirm, err := confirmStore.GetByCode(code)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		user, err := store.GetByEmail(confirm.Email)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		if user.Id != confirm.UserId {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("User Identification mismatch"))
			return
		}

		var json passwordReset
		err = c.BindJSON(&json)
		if err == nil {
			user.Password = json.Password
		} else {
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		err = store.PasswordSet(&user)
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

func PasswordResetInit(confirmStore confirmation.Store, store user.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error

		// TODO someone can spam this API endpoint and trigger lots of emails.

		var json passwordEmail
		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		user, err := store.GetByEmail(json.Email)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		// send confirmation email aka password reset
		confirm := confirmation.Record{
			Code:   util.RandSeq(20, util.LowerAlphaNumeric),
			UserId: user.Id,
			Email:  json.Email,
		}
		err = confirmStore.Create(&confirm)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		// TODO: Update language once we have an official company name
		// TODO: support mobile url
		u := url.Get(fmt.Sprintf("/password/reset/%s", confirm.Code))
		mailer := emailer.DefaultEmailer()
		err = mailer.SendMessage(
			json.Email,
			"noreply@cryptocades.com",
			"Password Reset",
			fmt.Sprintf("Hello! \nWe've recieved a request to reset your password. Please click the link below to continue. The link will expire after 12 hours.\n\n%s", u.String()),
		)
		if err != nil {
			log.Printf("Failed to send password reset: %s", err)
		}

		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	}
}
