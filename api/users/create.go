package users

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	recaptcha "github.com/ezzarghili/recaptcha-go"
	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/store/confirmation"
	"github.com/cbarraford/cryptocades-backend/store/income"
	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/cbarraford/cryptocades-backend/util"
	"github.com/cbarraford/cryptocades-backend/util/email"
	"github.com/cbarraford/cryptocades-backend/util/url"
)

const (
	SignUpBonus   = 5
	ReferralBonus = 10
	MaxReferrals  = 10
)

func Create(store user.Store, incomeStore income.Store, confirmStore confirmation.Store, captcha recaptcha.ReCAPTCHA) func(*gin.Context) {
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
			record.ReferralCode = json.ReferralCode
		} else {
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		if err := util.ValidateEmail(record.Email); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// verify captcha code
		success, err := captcha.Verify(json.CaptchaCode, c.ClientIP())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if !success {
			c.AbortWithError(http.StatusPreconditionFailed, fmt.Errorf("ReCAPTCHA failed"))
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

		err = NewUserBonus(txn, record, store, incomeStore)
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

func NewUserBonus(txn newrelic.Transaction, record user.Record, store user.Store, incomeStore income.Store) error {

	// give them free tickets for signing up
	in := income.Record{
		UserId:    record.Id,
		SessionId: "Sign up Bonus",
		Amount:    SignUpBonus,
	}
	seg := newrelic.DatastoreSegment{
		Product:    newrelic.DatastorePostgres,
		Collection: "incomes",
		Operation:  "sign-up-bonus",
	}
	seg.StartTime = newrelic.StartSegmentNow(txn)
	err := incomeStore.Create(&in)
	seg.End()
	if err != nil {
		return err
	}

	if record.ReferralCode != "" {

		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "incomes",
			Operation:  "Bonus Count",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		count, err := incomeStore.CountBonuses(record.Id, "Referral")
		seg.End()
		if err != nil {
			return err
		}

		if count < MaxReferrals {
			seg = newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: "users",
				Operation:  "GET",
			}
			seg.StartTime = newrelic.StartSegmentNow(txn)
			referrer, err := store.GetByReferralCode(record.ReferralCode)
			seg.End()
			if err != nil {
				return err
			}

			seg = newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: "users",
				Operation:  "GET",
			}
			seg.StartTime = newrelic.StartSegmentNow(txn)
			record, err = store.Get(record.Id)
			seg.End()
			if err != nil {
				return err
			}

			// give them free tickets for using referral
			in1 := income.Record{
				UserId:    record.Id,
				SessionId: fmt.Sprintf("Referral - %s", referrer.ReferralCode),
				Amount:    ReferralBonus,
			}
			in2 := income.Record{
				UserId:    referrer.Id,
				SessionId: fmt.Sprintf("Referral - %s", record.ReferralCode),
				Amount:    ReferralBonus,
			}

			seg = newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: "incomes",
				Operation:  "referral",
			}
			seg.StartTime = newrelic.StartSegmentNow(txn)
			err = incomeStore.Create(&in1)
			seg.End()
			if err != nil {
				return err
			}

			seg = newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: "incomes",
				Operation:  "referral",
			}
			seg.StartTime = newrelic.StartSegmentNow(txn)
			err = incomeStore.Create(&in2)
			seg.End()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
