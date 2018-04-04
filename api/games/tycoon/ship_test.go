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

func (s *ShipSuite) TestShipLogs(c *check.C) {
	gin.SetMode(gin.ReleaseMode)
	// happy path
	store := &mockStore{
		userId: 12,
	}

	r := gin.New()
	r.Use(middleware.TestSuite())
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.GET("/games/2/ships/:id/logs", GetShipLogs(store))

	req, _ := http.NewRequest("GET", "/games/2/ships/8/logs", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", "12")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var lines []asteroid_tycoon.Log
	c.Assert(json.Unmarshal(w.Body.Bytes(), &lines), check.IsNil)
	c.Assert(lines, check.HasLen, 1)
	line := lines[0]
	c.Check(line.Log, check.Equals, "log-line-text")
}

func (s *ShipSuite) TestShipUpgrade(c *check.C) {
	gin.SetMode(gin.ReleaseMode)
	// happy path
	store := &mockStore{
		userId: 12,
	}

	r := gin.New()
	r.Use(middleware.TestSuite())
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.PUT("/games/2/ships/:id/upgrade", ApplyUpgrade(store))
	body := strings.NewReader(`{"category_id": 4, "asset_id": 2}`)
	req, _ := http.NewRequest("PUT", "/games/2/ships/8/upgrade", body)
	req.Header.Set("Masquerade", "12")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)
	c.Check(store.updated, check.Equals, true)
	c.Check(store.assetId, check.Equals, 2)
	c.Check(store.categoryId, check.Equals, 4)
}

func (s *ShipSuite) TestMyAsteroids(c *check.C) {
	gin.SetMode(gin.ReleaseMode)
	store := &mockStore{
		userId: 12,
	}

	r := gin.New()
	r.Use(middleware.TestSuite())
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.GET("/games/2/ships/:id/asteroids", GetMyAsteroids(store))

	// happy path
	req, _ := http.NewRequest("GET", "/games/2/ships/8/asteroids", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", "12")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var asteroid asteroid_tycoon.Asteroid
	c.Assert(json.Unmarshal(w.Body.Bytes(), &asteroid), check.IsNil)
	c.Check(asteroid.ShipId, check.Equals, int64(8))
}

func (s *ShipSuite) TestGetStatus(c *check.C) {
	gin.SetMode(gin.ReleaseMode)
	store := &mockStore{
		userId: 12,
	}

	r := gin.New()
	r.Use(middleware.TestSuite())
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.GET("/games/2/ships/:id/status", GetStatus(store))

	// happy path
	req, _ := http.NewRequest("GET", "/games/2/ships/8/status", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", "12")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var status asteroid_tycoon.ShipStatus
	c.Assert(json.Unmarshal(w.Body.Bytes(), &status), check.IsNil)
	c.Check(status.Asteroid.ShipId, check.Equals, int64(8))
}
