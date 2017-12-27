package users

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"

	"github.com/CBarraford/lotto/api/middleware"
	"github.com/CBarraford/lotto/store/session"
)

type UserLogoutSuite struct{}

var _ = Suite(&UserLogoutSuite{})

type mockLogoutSessionStore struct {
	session.Dummy
	deleted bool
}

func (m *mockLogoutSessionStore) Delete(token string) error {
	m.deleted = true
	return nil
}

func (s *UserLogoutSuite) TestLogout(c *C) {
	gin.SetMode(gin.ReleaseMode)

	// happy path
	store := &mockLogoutSessionStore{}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.DELETE("/logout", Logout(store))
	req, _ := http.NewRequest("DELETE", "/logout", nil)
	req.Header.Set("Masquerade", "5")
	req.Header.Set("Session", "123456789")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200, Commentf("Response: %+v", w))
	c.Assert(store.deleted, Equals, true)
}
