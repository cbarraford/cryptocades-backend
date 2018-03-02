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
	"github.com/cbarraford/cryptocades-backend/store/income"
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

func (m *mockFacebookStore) GetByReferralCode(code string) (user.Record, error) {
	return user.Record{
		Id:           5,
		ReferralCode: code,
	}, nil
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

func (m *mockFacebookStore) Get(id int64) (user.Record, error) {
	return user.Record{
		Id:           id,
		Email:        "bob@bob.com",
		ReferralCode: fmt.Sprintf("ref-%d", id),
	}, nil
}

type mockIncomeStore struct {
	income.Dummy
	session_id []string
	freebe     bool
	count      int
	total      int
}

func (m *mockIncomeStore) Create(record *income.Record) error {
	m.session_id = append(m.session_id, record.SessionId)
	if record.Amount == 5 && record.SessionId == "Sign up Bonus" {
		m.freebe = true
	}
	m.total = m.total + record.Amount
	return nil
}

func (m *mockIncomeStore) CountBonuses(i int64, p string) (int, error) {
	return m.count, nil
}

func (s *FacebookLoginSuite) TestFacebookLogin(c *C) {
	gin.SetMode(gin.ReleaseMode)

	// happy path
	store := &mockFacebookStore{
		err: sql.ErrNoRows,
	}
	incomeStore := &mockIncomeStore{}
	sessionStore := &mockSessionStore{}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.Use(middleware.HandleErrors())
	r.POST("/login/facebook", Login(store, incomeStore, sessionStore))
	input := fmt.Sprintf(`{"email":"bob@bob.com","accessToken":"1234566789","referral_code":"code1"}`)
	body := strings.NewReader(input)
	req, _ := http.NewRequest("POST", "/login/facebook", body)
	req.Header.Set("Masquerade", "5")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200, Commentf("Response: %+v", w))
	c.Check(store.created, Equals, true)
	c.Check(store.user.Id, Equals, int64(12))

	c.Check(incomeStore.freebe, Equals, true)
	c.Check(incomeStore.session_id, DeepEquals, []string{"Sign up Bonus", "Referral - code1", "Referral - ref-12"})
	c.Check(incomeStore.total, Equals, 25) // 5 for bonus, 10 for each user (2)

}
