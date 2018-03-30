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

type AsteroidSuite struct{}

var _ = check.Suite(&AsteroidSuite{})

func (s *AsteroidSuite) TestAssignAsteroid(c *check.C) {
	gin.SetMode(gin.ReleaseMode)
	store := &mockStore{
		userId: 12,
	}

	r := gin.New()
	r.Use(middleware.TestSuite())
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.POST("/games/2/asteroids/assign", AssignAsteroid(store))

	// happy path
	body := strings.NewReader(`{"ship_id": 5, "asteroid_id": 4}`)
	req, _ := http.NewRequest("POST", "/games/2/asteroids/assign", body)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", "12")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	c.Check(store.created, check.Equals, true)
	c.Check(store.shipId, check.Equals, int64(5))
	c.Check(store.asteroidId, check.Equals, int64(4))
}

func (s *AsteroidSuite) TestAvailableAsteroids(c *check.C) {
	gin.SetMode(gin.ReleaseMode)
	store := &mockStore{
		userId: 12,
	}

	r := gin.New()
	r.Use(middleware.TestSuite())
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.GET("/games/2/asteroids/available", GetAvailableAsteroids(store))

	// happy path
	req, _ := http.NewRequest("GET", "/games/2/asteroids/available", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", "12")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var asteroids []asteroid_tycoon.Asteroid
	c.Assert(json.Unmarshal(w.Body.Bytes(), &asteroids), check.IsNil)
	c.Assert(asteroids, check.HasLen, 1)
	asteroid := asteroids[0]
	c.Check(asteroid.ShipId, check.Equals, int64(0))
}
