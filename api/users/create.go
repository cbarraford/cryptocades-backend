package users

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CBarraford/lotto/store/confirmation"
	"github.com/CBarraford/lotto/store/user"
	"github.com/CBarraford/lotto/util"
	"github.com/CBarraford/lotto/util/email"
)

func Create(store user.Store, confirmStore confirmation.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		record := user.Record{}
		var json input
		err = c.BindJSON(&json)
		if err == nil {
			record.Username = json.Username
			record.Email = json.Email
			record.BTCAddr = json.BTCAddr
			record.Password = json.Password
		} else {
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		err = store.Create(&record)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// send confirmation email
		confirm := confirmation.Record{
			Code:   util.RandSeq(20, util.LowerAlphaNumeric),
			UserId: record.Id,
			Email:  record.Email,
		}
		err = confirmStore.Create(&confirm)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		// TODO: Update language once we have an official company name
		// TODO: Update code URL once the frontend supports it
		emailer := email.DefaultEmailer()
		err = emailer.SendMessage(
			record.Email,
			"me@lotto.com",
			"Please confirm your email address", fmt.Sprintf("Hello! \nThanks for signing up for lotto. You must confirm your email address before you can start playing!\n\n%s", confirm.Code),
		)
		if err != nil {
			log.Printf("Failed to send email confirmation: %s", err)
		}

		c.JSON(http.StatusOK, record)
	}
}
