package users

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/api/middleware"
	"github.com/cbarraford/cryptocades-backend/store/session"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

type UserLoginSuite struct{}

var _ = Suite(&UserLoginSuite{})

type mockSessionStore struct {
	session.Dummy
	created bool
}

func (m *mockSessionStore) Create(record *session.Record, length int) error {
	m.created = true
	return nil
}

type mockUserStore struct {
	user.Dummy
}

func (m *mockUserStore) Authenticate(u, p string) (user.Record, error) {
	return user.Record{
		Id: 5,
	}, nil
}

func (s *UserLoginSuite) TestLogin(c *C) {
	gin.SetMode(gin.ReleaseMode)

	// happy path
	store := &mockUserStore{}
	sessionStore := &mockSessionStore{}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.POST("/login", Login(store, sessionStore))
	input := fmt.Sprintf(`{"username":"bob","password":"password"}`)
	body := strings.NewReader(input)
	req, _ := http.NewRequest("POST", "/login", body)
	req.Header.Set("Masquerade", "5")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200, Commentf("Response: %+v", w))
}
