package users

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	recaptcha "github.com/ezzarghili/recaptcha-go"
	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/api/middleware"
	"github.com/cbarraford/cryptocades-backend/store/boost"
	"github.com/cbarraford/cryptocades-backend/store/confirmation"
	"github.com/cbarraford/cryptocades-backend/store/income"
	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/cbarraford/cryptocades-backend/util/email"
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
	if record.Amount == 5 && record.SessionId == "Sign up Bonus" {
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

type mockBoostStore struct {
	boost.Dummy
	count int
}

func (m *mockBoostStore) Create(record *boost.Record) error {
	m.count = m.count + 1
	return nil
}

type mockReCAPTCHAClient struct{}

func (*mockReCAPTCHAClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	resp = &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
	}
	resp.Body = ioutil.NopCloser(strings.NewReader(`
    {
        "success": true,
        "challenge_ts": "2018-03-06T03:41:29+00:00",
        "hostname": "test.com"
    }
    `))
	return
}

func (s *UserCreateSuite) TestCreate(c *C) {
	gin.SetMode(gin.ReleaseMode)

	// happy path
	store := &mockCreateUserStore{}
	boostStore := &mockBoostStore{
		count: 0,
	}
	confirmStore := &mockConfirmCreateStore{}
	incomeStore := &mockIncomeCreateStore{
		count: 0,
	}
	captcha := recaptcha.ReCAPTCHA{
		Client: &mockReCAPTCHAClient{},
	}
	emailer, err := email.DefaultEmailer("../..")
	c.Assert(err, IsNil)

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.Use(middleware.HandleErrors())
	r.POST("/users", Create(store, incomeStore, confirmStore, boostStore, captcha, emailer))
	input := fmt.Sprintf(`{"username":"bob","password":"password","email":"bob@bob.com","btc_address":"12345","referrer":"code1"}`)
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
	c.Check(incomeStore.total, Equals, 5) // 5 for bonus
	c.Check(boostStore.count, Equals, 2)

	// make sure that when we have 10 referrals already, we don't award more
	incomeStore.count = 10
	incomeStore.session_id = nil
	r = gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.POST("/users", Create(store, incomeStore, confirmStore, boostStore, captcha, emailer))
	input = fmt.Sprintf(`{"username":"bob","password":"password","email":"bob@bob.com","btc_address":"12345","referrer":"code1"}`)
	body = strings.NewReader(input)
	req, _ = http.NewRequest("POST", "/users", body)
	req.Header.Set("Masquerade", "5")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)
	c.Check(incomeStore.session_id, DeepEquals, []string{"Sign up Bonus"})

	// test that a bad email fails
	incomeStore.count = 0
	incomeStore.session_id = nil
	r = gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.POST("/users", Create(store, incomeStore, confirmStore, boostStore, captcha, emailer))
	input = fmt.Sprintf(`{"username":"bob","password":"password","email":"bob+tag@bob.com","btc_address":"12345","referrer":"code1"}`)
	body = strings.NewReader(input)
	req, _ = http.NewRequest("POST", "/users", body)
	req.Header.Set("Masquerade", "5")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 400)

	// test that a bad username fails
	incomeStore.count = 0
	incomeStore.session_id = nil
	r = gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.POST("/users", Create(store, incomeStore, confirmStore, boostStore, captcha, emailer))
	input = fmt.Sprintf(`{"username":"bad username","password":"password","email":"bob@bobber.com","btc_address":"12345","referrer":"code1"}`)
	body = strings.NewReader(input)
	req, _ = http.NewRequest("POST", "/users", body)
	req.Header.Set("Masquerade", "5")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 400)
}
