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
	"github.com/cbarraford/cryptocades-backend/util/email"
	"github.com/cbarraford/cryptocades-backend/util/url"
)

func UpdateEmail(store user.Store, confirmStore confirmation.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		var userId int64
		var json passwordEmail
		var newEmail string

		err = c.BindJSON(&json)
		if err == nil {
			newEmail = json.Email
		} else {
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		if userId, err = context.GetUserId(c); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		record, err := store.Get(userId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// send confirmation email
		confirm := confirmation.Record{
			Code:   util.RandSeq(20, util.LowerAlphaNumeric),
			UserId: userId,
			Email:  newEmail,
		}
		err = confirmStore.Create(&confirm)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		// TODO: support mobile url
		u := url.Get(fmt.Sprintf("/confirmation/%s", confirm.Code))
		emailer := email.DefaultEmailer()
		err = emailer.SendMessage(
			record.Email,
			"noreply@cryptocades.com",
			"Please confirm your new email address",
			fmt.Sprintf("Hello! \nA change in email address is being registered with Cryptocades. You must confirm the new email address before it can take affect. Click the link below to confirm this new email address, %s\n\n%s", newEmail, u.String()),
		)
		if err != nil {
			log.Printf("Failed to send new email confirmation: %s", err)
		}

		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	}
}
