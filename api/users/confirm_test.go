package users

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store/confirmation"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

type UserConfirmSuite struct{}

var _ = Suite(&UserConfirmSuite{})

type mockConfirmStore struct {
	confirmation.Dummy
	email   string
	deleted bool
}

func (m *mockConfirmStore) GetByCode(code string) (confirmation.Record, error) {
	return confirmation.Record{
		Id:     1,
		Code:   code,
		Email:  m.email,
		UserId: 5,
	}, nil
}

func (m *mockConfirmStore) Delete(id int64) error {
	m.deleted = true
	return nil
}

type mockUserConfirmStore struct {
	user.Dummy
	confirmed bool
	email     string
}

func (*mockUserConfirmStore) Get(id int64) (user.Record, error) {
	return user.Record{
		Id:        id,
		Email:     "bob@bob.com",
		Confirmed: false,
	}, nil
}

func (m *mockUserConfirmStore) MarkAsConfirmed(record *user.Record) error {
	m.confirmed = true
	m.email = record.Email
	return nil
}

func (s *UserConfirmSuite) TestConfirm(c *C) {
	gin.SetMode(gin.ReleaseMode)

	// happy path
	confirmStore := &mockConfirmStore{email: "bob@cryptocades.com"}
	userStore := &mockUserConfirmStore{}

	r := gin.New()
	r.POST("/users/confirmation/:code", Confirm(confirmStore, userStore))
	req, _ := http.NewRequest("POST", "/users/confirmation/abcderf", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)
	c.Check(userStore.confirmed, Equals, true)
	c.Check(userStore.email, Equals, "bob@cryptocades.com")
	c.Check(confirmStore.deleted, Equals, true)
}
