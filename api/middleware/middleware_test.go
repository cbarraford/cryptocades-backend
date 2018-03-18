package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cbarraford/cryptocades-backend/store/session"
	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type MiddlewareSuite struct{}

var _ = Suite(&MiddlewareSuite{})

func (s *MiddlewareSuite) TestMasquerade(c *C) {
	r := gin.New()
	r.Use(Masquerade())
	r.Use(AuthRequired())

	masqueradeId := "myId"
	r.GET("/test", func(context *gin.Context) {
		userId, _ := context.Get("userId")
		c.Check(userId, Equals, masqueradeId, Commentf("UserId: %s", userId))
		escalated, ok := context.Get("escalated")
		c.Assert(ok, Equals, true)
		c.Check(escalated, Equals, true)
		context.String(200, "OK")
	})

	// happy path
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", masqueradeId)
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)

	// no auth in request
	req, _ = http.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 401)
}

type mockSessionStore struct {
	session.Dummy
	authed    bool
	escalated bool
	admin     bool
}

func (m *mockSessionStore) Authenticate(token string) (int64, bool, bool, error) {
	m.authed = true
	return 5, m.escalated, m.admin, nil
}

func (s *MiddlewareSuite) TestEscalatedAuthenticate(c *C) {
	sessionStore := &mockSessionStore{escalated: true}

	r := gin.New()
	r.Use(Authenticate(sessionStore))
	r.Use(EscalatedAuthRequired())

	r.GET("/test", func(context *gin.Context) {
		userId, _ := context.Get("userId")
		c.Check(userId, Equals, "5", Commentf("UserId: %s", userId))
		context.String(200, "OK")
	})

	// happy path
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Session", "a session id is here")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)

	// no auth in request
	req, _ = http.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 401)

	sessionStore = &mockSessionStore{escalated: false}

	r = gin.New()
	r.Use(Authenticate(sessionStore))
	r.Use(EscalatedAuthRequired())

	r.GET("/test", func(context *gin.Context) {
		userId, _ := context.Get("userId")
		c.Check(userId, Equals, "5", Commentf("UserId: %s", userId))
		context.String(200, "OK")
	})

	req, _ = http.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	req.Header.Set("Session", "a session id is here")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 401)

}

func (s *MiddlewareSuite) TestAuthenticate(c *C) {
	sessionStore := &mockSessionStore{}

	r := gin.New()
	r.Use(Authenticate(sessionStore))
	r.Use(AuthRequired())

	r.GET("/test", func(context *gin.Context) {
		userId, _ := context.Get("userId")
		c.Check(userId, Equals, "5", Commentf("UserId: %s", userId))
		context.String(200, "OK")
	})

	// happy path
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Session", "a session id is here")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)

	// no auth in request
	req, _ = http.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 401)
}

func (s *MiddlewareSuite) TestAuthAuthentication(c *C) {
	sessionStore := &mockSessionStore{
		admin: true,
	}

	r := gin.New()
	r.Use(Authenticate(sessionStore))
	r.Use(AdminAuthRequired())

	r.GET("/test", func(context *gin.Context) {
		userId, _ := context.Get("userId")
		c.Check(userId, Equals, "5", Commentf("UserId: %s", userId))
		context.String(200, "OK")
	})
	r.GET("/test/bot", func(context *gin.Context) {
		context.String(200, "OK")
	})

	// happy path
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Session", "a session id is here")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)

	req, _ = http.NewRequest("GET", "/test/bot", nil)
	w = httptest.NewRecorder()
	req.Header.Set("Token", "QieDpVTtcnBgFVDPccRmDa98")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)

	// non-admin
	sessionStore.admin = false
	req, _ = http.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	req.Header.Set("Session", "a session id is here")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 401)
}
