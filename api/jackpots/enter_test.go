package jackpots

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/api/middleware"
	"github.com/cbarraford/cryptocades-backend/store/entry"
	"github.com/cbarraford/cryptocades-backend/store/user"
)

type JackpotEnterSuite struct{}

var _ = check.Suite(&JackpotEnterSuite{})

type mockEntryStore struct {
	entry.Dummy
	created   bool
	amount    int
	jackpotId int64
	userId    int64
}

func (*mockEntryStore) UserSpent(i int64) (int, error) {
	return 40, nil
}

func (s *mockEntryStore) Create(record *entry.Record) error {
	s.created = true
	s.amount = record.Amount
	s.jackpotId = record.JackpotId
	s.userId = record.UserId
	return nil
}

type mockUserStore struct {
	user.Dummy
}

func (*mockUserStore) Get(i int64) (user.Record, error) {
	return user.Record{
		Id:          i,
		MinedHashes: 50,
		BonusHashes: 5,
	}, nil
}

func (s *JackpotEnterSuite) TestEnter(c *check.C) {
	store := &mockEntryStore{}
	userStore := &mockUserStore{}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.POST("/jackpots/:id/enter", Enter(userStore, store))

	// happy path
	input := fmt.Sprintf(`{"amount":10}`)
	body := strings.NewReader(input)
	req, _ := http.NewRequest("POST", "/jackpots/23/enter", body)
	req.Header.Set("Masquerade", "44")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)
	c.Check(store.created, check.Equals, true)
	c.Check(store.amount, check.Equals, 10)
	c.Check(store.jackpotId, check.Equals, int64(23))
	c.Check(store.userId, check.Equals, int64(44))

	// overspend
	store.created = false
	input = fmt.Sprintf(`{"amount":1000}`)
	body = strings.NewReader(input)
	req, _ = http.NewRequest("POST", "/jackpots/23/enter", body)
	req.Header.Set("Masquerade", "44")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 402)
	c.Check(store.created, check.Equals, false)
}
