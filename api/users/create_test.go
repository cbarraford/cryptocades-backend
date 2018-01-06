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

type UserCreateSuite struct{}

var _ = Suite(&UserCreateSuite{})

type mockCreateUserStore struct {
	user.Dummy
	btc      string
	username string
	password string
	email    string
	created  bool
}

func (m *mockCreateUserStore) Create(record *user.Record) error {
	record.Id = 10
	m.created = true
	m.btc = record.BTCAddr
	m.username = record.Username
	m.password = record.Password
	m.email = record.Email
	return nil
}

type mockConfirmCreateStore struct {
	confirmation.Dummy
	created bool
	code    string
	userId  int64
	email   string
}

func (m *mockConfirmCreateStore) Create(record *confirmation.Record) error {
	m.created = true
	m.userId = record.UserId
	m.email = record.Email
	m.code = record.Code
	return nil
}

func (s *UserCreateSuite) TestCreate(c *C) {
	gin.SetMode(gin.ReleaseMode)

	// happy path
	store := &mockCreateUserStore{}
	confirmStore := &mockConfirmCreateStore{}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.POST("/users", Create(store, confirmStore))
	input := fmt.Sprintf(`{"username":"bob","password":"password","email":"bob@bob.com","btc_address":"12345"}`)
	body := strings.NewReader(input)
	req, _ := http.NewRequest("POST", "/users", body)
	req.Header.Set("Masquerade", "5")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)
	c.Check(store.created, Equals, true)
	c.Check(store.btc, Equals, "12345")
	c.Check(store.username, Equals, "bob")
	c.Check(store.password, Equals, "password")
	c.Check(store.email, Equals, "bob@bob.com")

	c.Assert(confirmStore.created, Equals, true)
	c.Check(confirmStore.email, Equals, "bob@bob.com")
	c.Check(confirmStore.userId, Equals, int64(10))
	c.Check(confirmStore.code, Not(Equals), "")
}
