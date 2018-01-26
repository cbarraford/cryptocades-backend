package users

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store/confirmation"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

type UserPasswordResetSuite struct{}

var _ = Suite(&UserPasswordResetSuite{})

type mockPasswordResetStore struct {
	confirmation.Dummy
	email   string
	deleted bool
	created bool
}

func (m *mockPasswordResetStore) Create(r *confirmation.Record) error {
	m.created = true
	return nil
}

func (m *mockPasswordResetStore) GetByCode(code string) (confirmation.Record, error) {
	return confirmation.Record{
		Id:     1,
		Code:   code,
		Email:  m.email,
		UserId: 5,
	}, nil
}

func (m *mockPasswordResetStore) Delete(id int64) error {
	m.deleted = true
	return nil
}

type mockUserPasswordResetStore struct {
	user.Dummy
	reset    bool
	password string
}

func (*mockUserPasswordResetStore) GetByEmail(email string) (user.Record, error) {
	return user.Record{
		Id:    5,
		Email: email,
	}, nil
}

func (m *mockUserPasswordResetStore) PasswordSet(record *user.Record) error {
	m.reset = true
	m.password = record.Password
	return nil
}

func (s *UserPasswordResetSuite) TestPasswordReset(c *C) {
	gin.SetMode(gin.ReleaseMode)

	// happy path
	confirmStore := &mockPasswordResetStore{email: "bob@cryptocades.com"}
	userStore := &mockUserPasswordResetStore{}

	r := gin.New()
	r.POST("/users/password_reset", PasswordResetInit(confirmStore, userStore))
	r.POST("/users/password_reset/:code", PasswordReset(confirmStore, userStore))

	input := fmt.Sprintf(`{"email":"bob@cryptocades.com"}`)
	body := strings.NewReader(input)
	req, _ := http.NewRequest("POST", "/users/password_reset", body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)
	c.Check(confirmStore.created, Equals, true)

	input = fmt.Sprintf(`{"password":"new_password"}`)
	body = strings.NewReader(input)
	req, _ = http.NewRequest("POST", "/users/password_reset/abcderf", body)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)
	c.Check(userStore.reset, Equals, true)
	c.Check(userStore.password, Equals, "new_password")
	c.Check(confirmStore.deleted, Equals, true)
}
