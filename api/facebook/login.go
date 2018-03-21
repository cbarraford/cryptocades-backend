package facebook

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/api/users"
	"github.com/cbarraford/cryptocades-backend/store/boost"
	"github.com/cbarraford/cryptocades-backend/store/income"
	"github.com/cbarraford/cryptocades-backend/store/session"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

// TODO: don't get email from request, could be forged
type login struct {
	Email       string `json:"email"`
	ExpiresIn   int    `json:"expiresIn"`
	AccessToken string `json:"accessToken"`
	Referrer    string `json:"referrer"`
}

type facebookMe struct {
	Name   string `json:"name"`
	UserId string `json:"id"`
}

func Login(store user.Store, boostStore boost.Store, incomeStore income.Store, sessionStore session.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		var record user.Record
		var fb facebookMe
		var _login login

		txn := nrgin.Transaction(c)

		err = c.BindJSON(&_login)
		if err == nil {

			client := &http.Client{}
			q := fmt.Sprintf("https://graph.facebook.com/me?access_token=%s", _login.AccessToken)

			req, err := http.NewRequest("GET", q, nil)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			segEx := newrelic.StartExternalSegment(txn, req)
			resp, err := client.Do(req)
			segEx.Response = resp
			segEx.End()
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			err = json.NewDecoder(resp.Body).Decode(&fb)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			seg := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: "users",
				Operation:  "GET",
			}
			seg.StartTime = newrelic.StartSegmentNow(txn)
			record, err = store.GetByFacebookId(fb.UserId)
			seg.End()
			if err != nil {
				if err == sql.ErrNoRows {
					// user doesn't exist, create them.
					record.Email = _login.Email
					record.Username = fb.UserId
					record.FacebookId = fb.UserId
					record.Referrer = _login.Referrer
					// TODO Require ReCAPTCHA
					err := store.Create(&record)
					if err != nil {
						c.AbortWithError(http.StatusBadRequest, err)
						return
					}

					err = users.NewUserBonus(txn, record, store, incomeStore, boostStore)
					if err != nil {
						c.AbortWithError(http.StatusBadRequest, err)
						return
					}
				} else {
					c.AbortWithError(http.StatusBadRequest, err)
					return
				}
			}
		} else {
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		sessionRecord := session.Record{
			UserId: record.Id,
		}
		seg := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: "session",
			Operation:  "Create",
		}
		seg.StartTime = newrelic.StartSegmentNow(txn)
		err = sessionStore.Create(&sessionRecord, 30)
		seg.End()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, sessionRecord)
	}

}
