package jackpots

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/api/middleware"
	"github.com/cbarraford/cryptocades-backend/store/entry"
	"github.com/cbarraford/cryptocades-backend/store/jackpot"
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
	btcAddr string
}

func (s *mockUserStore) Get(id int64) (user.Record, error) {
	return user.Record{
		Id:      id,
		BTCAddr: s.btcAddr,
	}, nil
}

type mockJackpotStore struct {
	jackpot.Dummy
	endtime time.Time
}

func (s *mockJackpotStore) Get(id int64) (jackpot.Record, error) {
	return jackpot.Record{
		Id:      id,
		EndTime: s.endtime,
	}, nil
}

func (s *JackpotEnterSuite) TestEnter(c *check.C) {
	store := &mockEntryStore{}
	userStore := &mockUserStore{
		btcAddr: "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq",
	}
	jackpotStore := &mockJackpotStore{
		endtime: time.Now().Add(168 * time.Hour),
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.POST("/jackpots/:id/enter", Enter(store, userStore, jackpotStore))

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

	// invalid btc address
	userStore.btcAddr = "           "
	req, _ = http.NewRequest("POST", "/jackpots/23/enter", body)
	req.Header.Set("Masquerade", "44")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 400)

	// only enter active jackpots
	userStore.btcAddr = "1MiJFQvupX5kSZcUtfSoD9NtLevUgjv3uq"
	jackpotStore.endtime = time.Now().Add(-1 * time.Second)
	req, _ = http.NewRequest("POST", "/jackpots/23/enter", body)
	req.Header.Set("Masquerade", "44")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 400)

}
