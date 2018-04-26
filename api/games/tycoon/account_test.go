package tycoon

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/api/middleware"
	"github.com/cbarraford/cryptocades-backend/store/game/asteroid_tycoon"
)

type AccountSuite struct{}

var _ = check.Suite(&AccountSuite{})

func (s *AccountSuite) TestCreateAccount(c *check.C) {
	gin.SetMode(gin.ReleaseMode)
	store := &mockStore{}

	r := gin.New()
	r.Use(middleware.HandleErrors())
	r.Use(middleware.Masquerade())
	r.POST("/games/2/account", CreateAccount(store))

	// happy path
	req, _ := http.NewRequest("POST", "/games/2/account", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", "12")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var acct asteroid_tycoon.Account
	c.Assert(json.Unmarshal(w.Body.Bytes(), &acct), check.IsNil)
	c.Check(acct.UserId, check.Equals, int64(12))
	c.Check(store.created, check.Equals, true)
	c.Check(store.userId, check.Equals, int64(12))
}

func (s *AccountSuite) TestGetAccount(c *check.C) {
	gin.SetMode(gin.ReleaseMode)
	store := &mockStore{}

	r := gin.New()
	r.Use(middleware.HandleErrors())
	r.Use(middleware.Masquerade())
	r.GET("/games/2/account", GetAccount(store))

	// happy path
	req, _ := http.NewRequest("GET", "/games/2/account", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", "12")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var acct asteroid_tycoon.Account
	c.Assert(json.Unmarshal(w.Body.Bytes(), &acct), check.IsNil)
	c.Check(acct.Id, check.Equals, int64(1))
	c.Check(acct.UserId, check.Equals, int64(12))
}

func (s *AccountSuite) TestTradeForCredits(c *check.C) {
	gin.SetMode(gin.ReleaseMode)
	store := &mockStore{}

	r := gin.New()
	r.Use(middleware.TestSuite())
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.POST("/games/2/trade/credits", TradeForCredits(store))

	// happy path
	body := strings.NewReader(`{"amount":5}`)
	req, _ := http.NewRequest("POST", "/games/2/trade/credits", body)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", "12")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var ship asteroid_tycoon.Ship
	c.Assert(json.Unmarshal(w.Body.Bytes(), &ship), check.IsNil)
	c.Check(store.accountId, check.Equals, int64(1))
	c.Check(store.amount, check.Equals, 5)
}

func (s *AccountSuite) TestTradeForPlays(c *check.C) {
	gin.SetMode(gin.ReleaseMode)
	store := &mockStore{}

	r := gin.New()
	r.Use(middleware.TestSuite())
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.POST("/games/2/trade/plays", TradeForPlays(store))

	// happy path
	body := strings.NewReader(`{"amount":5}`)
	req, _ := http.NewRequest("POST", "/games/2/trade/plays", body)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", "12")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var ship asteroid_tycoon.Ship
	c.Assert(json.Unmarshal(w.Body.Bytes(), &ship), check.IsNil)
	c.Check(store.accountId, check.Equals, int64(1))
	c.Check(store.amount, check.Equals, 5)
}
