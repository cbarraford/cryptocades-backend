package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

func TestPackage(t *testing.T) { TestingT(t) }

type ApiSuite struct{}

var _ = Suite(&ApiSuite{})

func (s *ApiSuite) TestApiService(c *C) {
	store := store.Store{
		Users: &user.Dummy{},
	}

	r := GetAPIService(store)

	// check ping apiendpoint
	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)
}
