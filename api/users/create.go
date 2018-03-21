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

	"github.com/cbarraford/cryptocades-backend/store/boost"
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

func Create(store user.Store, incomeStore income.Store, confirmStore confirmation.Store, boostStore boost.Store, captcha recaptcha.ReCAPTCHA, emailer email.Emailer) func(*gin.Context) {
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
			record.Referrer = json.Referrer
		} else {
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		if err := util.ValidateEmail(record.Email); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if err := util.ValidateUsername(record.Username); err != nil {
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

		err = NewUserBonus(txn, record, store, incomeStore, boostStore)
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

		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "users",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		record, err = store.Get(record.Id)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// TODO: support mobile url
		segEmail := newrelic.StartSegment(txn, "email")
		emailTemplate := email.EmailTemplate{
			Subject:     "Confirm your Cryptocades account",
			ConfirmURL:  url.Get(fmt.Sprintf("/confirmation/%s", confirm.Code)).String(),
			ReferralURL: url.Get(fmt.Sprintf("/signup?referral=%s", record.Referrer)).String(),
		}
		err = emailer.SendHTML(
			record.Email,
			"noreply@cryptocades.com",
			emailTemplate.Subject,
			"confirm",
			emailTemplate,
		)
		segEmail.End()
		if err != nil {
			log.Printf("Failed to send email confirmation: %s", err)
		}

		c.JSON(http.StatusOK, record)
	}
}

func NewUserBonus(txn newrelic.Transaction, record user.Record, store user.Store, incomeStore income.Store, boostStore boost.Store) error {

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

	if record.Referrer != "" {

		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "users",
			Operation:  "GET",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		referrer, err := store.GetByReferralCode(record.Referrer)
		seg.End()
		if err != nil {
			return err
		}

		// grant referrer a boost
		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "boosts",
			Operation:  "CREATE",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = boostStore.Create(&boost.Record{UserId: referrer.Id})
		seg.End()
		if err != nil {
			return err
		}

		// grant new user a boost
		seg = newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "boosts",
			Operation:  "CREATE",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = boostStore.Create(&boost.Record{UserId: record.Id})
		seg.End()
		if err != nil {
			return err
		}
	}
	return nil
}
