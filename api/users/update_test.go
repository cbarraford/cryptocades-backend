package users

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"

	"github.com/CBarraford/lotto/api/middleware"
	"github.com/CBarraford/lotto/store/user"
)

type UserUpdateSuite struct{}

var _ = Suite(&UserUpdateSuite{})

type mockUpdateStore struct {
	user.Dummy
	updated     bool
	username    string
	btc_address string
	email       string
}

func (*mockUpdateStore) Get(id int64) (user.Record, error) {
	return user.Record{
		Id: id,
	}, nil
}

func (m *mockUpdateStore) Update(record *user.Record) error {
	m.updated = true
	m.username = record.Username
	m.email = record.Email
	m.btc_address = record.BTCAddr
	return nil
}

func (s *UserUpdateSuite) TestUpdate(c *C) {
	gin.SetMode(gin.ReleaseMode)

	// happy path
	store := &mockUpdateStore{}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.PUT("/me", Update(store))
	body := strings.NewReader(`{"email":"chad@test.com","btc_address":"1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3ub","username":"jonny"}`)
	req, _ := http.NewRequest("PUT", "/me", body)
	req.Header.Set("Masquerade", "5")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)
	c.Check(store.updated, Equals, true)
	c.Check(store.username, Equals, "jonny")
	c.Check(store.email, Equals, "chad@test.com")
	c.Check(store.btc_address, Equals, "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3ub")
}
