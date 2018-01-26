package users

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/api/middleware"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

type UserDeleteSuite struct{}

var _ = Suite(&UserDeleteSuite{})

type mockDeleteStore struct {
	user.Dummy
	deleted bool
	id      int64
}

func (m *mockDeleteStore) Delete(id int64) error {
	m.deleted = true
	m.id = id
	return nil
}

func (s *UserDeleteSuite) TestDelete(c *C) {
	gin.SetMode(gin.ReleaseMode)

	// happy path
	store := &mockDeleteStore{}
	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.DELETE("/me", Delete(store))
	req, _ := http.NewRequest("DELETE", "/me", nil)
	req.Header.Set("Masquerade", "5")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)
	c.Check(store.deleted, Equals, true)
	c.Check(store.id, Equals, int64(5))
}
