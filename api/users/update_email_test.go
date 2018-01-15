package users

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"

	"github.com/CBarraford/lotto/api/middleware"
	"github.com/CBarraford/lotto/store/confirmation"
	"github.com/CBarraford/lotto/store/user"
)

type UserUpdateEmailSuite struct{}

var _ = Suite(&UserUpdateEmailSuite{})

type mockUpdateEmailUserStore struct {
	user.Dummy
}

func (m *mockUpdateEmailUserStore) Get(id int64) (user.Record, error) {
	return user.Record{
		Id:    id,
		Email: "bobby@bobs.com",
	}, nil
}

type mockConfirmUpdateEmailStore struct {
	confirmation.Dummy
	created bool
	code    string
	userId  int64
	email   string
}

func (m *mockConfirmUpdateEmailStore) Create(record *confirmation.Record) error {
	m.created = true
	m.userId = record.UserId
	m.email = record.Email
	m.code = record.Code
	return nil
}

func (s *UserUpdateEmailSuite) TestUpdateEmail(c *C) {
	gin.SetMode(gin.ReleaseMode)

	// happy path
	store := &mockUpdateEmailUserStore{}
	confirmStore := &mockConfirmUpdateEmailStore{}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.EscalatedAuthRequired())
	r.PUT("/me/email", UpdateEmail(store, confirmStore))
	input := fmt.Sprintf(`{"email":"bob@bob.com"}`)
	body := strings.NewReader(input)
	req, _ := http.NewRequest("PUT", "/me/email", body)
	req.Header.Set("Masquerade", "5")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)

	c.Assert(confirmStore.created, Equals, true)
	c.Check(confirmStore.email, Equals, "bob@bob.com")
	c.Check(confirmStore.userId, Equals, int64(5))
	c.Check(confirmStore.code, Not(Equals), "")
}
