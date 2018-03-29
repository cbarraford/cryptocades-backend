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

type ShipSuite struct{}

var _ = check.Suite(&ShipSuite{})

func (s *ShipSuite) TestCreateShip(c *check.C) {
	gin.SetMode(gin.ReleaseMode)
	store := &mockStore{}

	r := gin.New()
	r.Use(middleware.TestSuite())
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.POST("/games/2/ships", CreateShip(store))

	// happy path
	req, _ := http.NewRequest("POST", "/games/2/ships", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", "12")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var ship asteroid_tycoon.Ship
	c.Assert(json.Unmarshal(w.Body.Bytes(), &ship), check.IsNil)
	c.Check(ship.AccountId, check.Equals, int64(1))
}

func (s *ShipSuite) TestGetShips(c *check.C) {
	gin.SetMode(gin.ReleaseMode)
	store := &mockStore{
		userId: 12,
	}

	r := gin.New()
	r.Use(middleware.TestSuite())
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.GET("/games/2/ships", GetShips(store))

	// happy path
	req, _ := http.NewRequest("GET", "/games/2/ships", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", "12")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var ships []asteroid_tycoon.Ship
	c.Assert(json.Unmarshal(w.Body.Bytes(), &ships), check.IsNil)
	c.Assert(ships, check.HasLen, 1)
	ship := ships[0]
	c.Check(ship.AccountId, check.Equals, int64(1))
}

func (s *ShipSuite) TestShipUpdate(c *check.C) {
	gin.SetMode(gin.ReleaseMode)
	// happy path
	store := &mockStore{
		userId: 12,
	}

	r := gin.New()
	r.Use(middleware.TestSuite())
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.PUT("/games/2/ships/:id", UpdateShip(store))
	body := strings.NewReader(`{"name":"bulldozer"}`)
	req, _ := http.NewRequest("PUT", "/games/2/ships/8", body)
	req.Header.Set("Masquerade", "12")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)
	c.Check(store.updated, check.Equals, true)
	c.Check(store.name, check.Equals, "bulldozer")
}
