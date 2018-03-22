package users

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/api/middleware"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

type UserUpdateSuite struct{}

var _ = Suite(&UserUpdateSuite{})

type mockUpdateStore struct {
	user.Dummy
	updated     bool
	btc_address string
	password    string
	username    string
}

func (*mockUpdateStore) Get(id int64) (user.Record, error) {
	return user.Record{
		Id:       id,
		Password: "oldPassword",
	}, nil
}

func (m *mockUpdateStore) Update(record *user.Record) error {
	m.updated = true
	m.btc_address = record.BTCAddr
	m.username = record.Username
	return nil
}

func (m *mockUpdateStore) PasswordSet(record *user.Record) error {
	m.password = record.Password
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
	body := strings.NewReader(`{"btc_address":"1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3ub","password":"newPassword", "username":"bobby"}`)
	req, _ := http.NewRequest("PUT", "/me", body)
	req.Header.Set("Masquerade", "5")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)
	c.Check(store.updated, Equals, true)
	c.Check(store.btc_address, Equals, "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3ub")
	c.Check(store.password, Equals, "newPassword")
	c.Check(store.username, Equals, "bobby")

	// bad btc address
	body = strings.NewReader(`{"btc_address":"bad btc address"}`)
	req, _ = http.NewRequest("PUT", "/me", body)
	req.Header.Set("Masquerade", "5")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 400)
}
