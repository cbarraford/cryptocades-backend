package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"

	"github.com/CBarraford/lotto/api/middleware"
	"github.com/CBarraford/lotto/store/entry"
	"github.com/CBarraford/lotto/store/user"
)

type UserBalanceSuite struct{}

var _ = check.Suite(&UserBalanceSuite{})

type mockBalanceEntriesStore struct {
	entry.Dummy
}

func (*mockBalanceEntriesStore) UserSpent(id int64) (int, error) {
	return 5, nil
}

type mockBalanceStore struct {
	user.Dummy
}

func (*mockBalanceStore) Get(id int64) (user.Record, error) {
	return user.Record{
		Id:          id,
		Username:    "bob",
		MinedHashes: 12,
		BonusHashes: 10,
	}, nil
}

type response struct {
	Balance int `json:"balance"`
}

func (s *UserBalanceSuite) TestBalance(c *check.C) {
	store := &mockBalanceStore{}
	entryStore := &mockBalanceEntriesStore{}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())

	r.GET("/me/balance", Balance(store, entryStore))

	// happy path
	req, _ := http.NewRequest("GET", "/me/balance", nil)
	req.Header.Set("Masquerade", "12")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	record := response{}
	c.Assert(json.Unmarshal(w.Body.Bytes(), &record), check.IsNil)
	c.Check(record.Balance, check.Equals, 17)
}
