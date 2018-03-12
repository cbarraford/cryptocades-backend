package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store"
	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/cbarraford/cryptocades-backend/util/email"
	recaptcha "github.com/ezzarghili/recaptcha-go"
	newrelic "github.com/newrelic/go-agent"
)

func TestPackage(t *testing.T) { TestingT(t) }

type ApiSuite struct{}

var _ = Suite(&ApiSuite{})

func (s *ApiSuite) TestApiService(c *C) {
	store := store.Store{
		Users: &user.Dummy{},
	}

	config := newrelic.NewConfig("Test", "")
	config.Enabled = false
	agent, err := newrelic.NewApplication(config)
	c.Assert(err, IsNil)

	r := GetAPIService(store, agent, recaptcha.ReCAPTCHA{}, email.Emailer{})

	// check ping apiendpoint
	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)
}
