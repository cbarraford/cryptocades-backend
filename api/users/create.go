package users

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/store/confirmation"
	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/cbarraford/cryptocades-backend/util"
	"github.com/cbarraford/cryptocades-backend/util/email"
	"github.com/cbarraford/cryptocades-backend/util/url"
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

		txn := nrgin.Transaction(c)
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "users",
			Operation:  "Create",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = store.Create(&record)
		seg.End()
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

		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "confirmations",
			Operation:  "Create",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = confirmStore.Create(&confirm)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// TODO: support mobile url
		u := url.Get(fmt.Sprintf("/confirmation/%s", confirm.Code))
		emailer := email.DefaultEmailer()
		segEmail := newrelic.StartSegment(txn, "email")
		err = emailer.SendMessage(
			record.Email,
			"noreply@cryptocades.com",
			"Please confirm your email address",
			fmt.Sprintf("Hello! \nThanks for signing up for Cryptocades. You must confirm your email address before you can start playing!\n\n%s", u.String()),
		)
		segEmail.End()
		if err != nil {
			log.Printf("Failed to send email confirmation: %s", err)
		}

		c.JSON(http.StatusOK, record)
	}
}
