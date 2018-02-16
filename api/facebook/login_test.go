package facebook

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/api/middleware"
	"github.com/cbarraford/cryptocades-backend/store/session"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

func TestPackage(t *testing.T) { TestingT(t) }

type FacebookLoginSuite struct{}

var _ = Suite(&FacebookLoginSuite{})

type mockSessionStore struct {
	session.Dummy
	created bool
}

func (m *mockSessionStore) Create(record *session.Record, length int) error {
	m.created = true
	return nil
}

type mockFacebookStore struct {
	user.Dummy
	created bool
	user    user.Record
	err     error
}

func (m *mockFacebookStore) GetByFacebookId(id string) (user.Record, error) {
	return m.user, m.err
}

func (m *mockFacebookStore) Create(record *user.Record) error {
	m.created = true
	record.Id = int64(12)
	m.user = *record
	return nil
}

func (s *FacebookLoginSuite) TestFacebookLogin(c *C) {
	gin.SetMode(gin.ReleaseMode)

	// happy path
	store := &mockFacebookStore{
		err: sql.ErrNoRows,
	}
	sessionStore := &mockSessionStore{}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.Use(middleware.HandleErrors())
	r.POST("/login/facebook", Login(store, sessionStore))
	input := fmt.Sprintf(`{"email":"bob@bob.com","accessToken":"1234566789"}`)
	body := strings.NewReader(input)
	req, _ := http.NewRequest("POST", "/login/facebook", body)
	req.Header.Set("Masquerade", "5")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200, Commentf("Response: %+v", w))
	c.Check(store.created, Equals, true)
	c.Check(store.user.Id, Equals, int64(12))
}
