package users

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/api/middleware"
	"github.com/cbarraford/cryptocades-backend/store/confirmation"
	"github.com/cbarraford/cryptocades-backend/store/income"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

type UserCreateSuite struct{}

var _ = Suite(&UserCreateSuite{})

type mockIncomeCreateStore struct {
	income.Dummy
	session_id []string
	freebe     bool
	count      int
	total      int
}

func (m *mockIncomeCreateStore) Create(record *income.Record) error {
	m.session_id = append(m.session_id, record.SessionId)
	if record.Amount == 5 && record.SessionId == "Sign up" {
		m.freebe = true
	}
	m.total = m.total + record.Amount
	return nil
}

func (m *mockIncomeCreateStore) CountBonuses(i int64, p string) (int, error) {
	return m.count, nil
}

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

func (m *mockCreateUserStore) GetByReferralCode(code string) (user.Record, error) {
	return user.Record{
		Id:           5,
		ReferralCode: code,
	}, nil
}

func (m *mockCreateUserStore) Get(id int64) (user.Record, error) {
	return user.Record{
		Id:           id,
		Email:        "bob@bob.com",
		ReferralCode: fmt.Sprintf("ref-%d", id),
	}, nil
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
	incomeStore := &mockIncomeCreateStore{
		count: 0,
	}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.POST("/users", Create(store, incomeStore, confirmStore))
	input := fmt.Sprintf(`{"username":"bob","password":"password","email":"bob@bob.com","btc_address":"12345","referral_code":"code1"}`)
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

	c.Check(incomeStore.freebe, Equals, true)
	c.Check(incomeStore.session_id, DeepEquals, []string{"Sign up", "Referral - code1", "Referral - ref-10"})
	c.Check(incomeStore.total, Equals, 25) // 5 for bonus, 10 for each user (2)

	// make sure that when we have 10 referrals already, we don't award more
	incomeStore.count = 10
	incomeStore.session_id = nil
	r = gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.POST("/users", Create(store, incomeStore, confirmStore))
	input = fmt.Sprintf(`{"username":"bob","password":"password","email":"bob@bob.com","btc_address":"12345","referral_code":"code1"}`)
	body = strings.NewReader(input)
	req, _ = http.NewRequest("POST", "/users", body)
	req.Header.Set("Masquerade", "5")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)
	c.Check(incomeStore.session_id, DeepEquals, []string{"Sign up"})

	// test that a bad email fails
	incomeStore.count = 0
	incomeStore.session_id = nil
	r = gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.POST("/users", Create(store, incomeStore, confirmStore))
	input = fmt.Sprintf(`{"username":"bob","password":"password","email":"bob+tag@bob.com","btc_address":"12345","referral_code":"code1"}`)
	body = strings.NewReader(input)
	req, _ = http.NewRequest("POST", "/users", body)
	req.Header.Set("Masquerade", "5")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 400)
}
