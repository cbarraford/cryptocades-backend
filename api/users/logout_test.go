package users

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

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
	input := fmt.Sprintf(`{"token":"12345"}`)
	body := strings.NewReader(input)
	req, _ := http.NewRequest("DELETE", "/logout", body)
	req.Header.Set("Masquerade", "5")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200, Commentf("Response: %+v", w))
	c.Assert(store.deleted, Equals, true)
}
