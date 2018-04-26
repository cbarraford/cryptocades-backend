package tycoon

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/api/middleware"
)

type LedgerSuite struct{}

var _ = check.Suite(&LedgerSuite{})

func (s *LedgerSuite) TestAssignAsteroid(c *check.C) {
	gin.SetMode(gin.ReleaseMode)
	store := &mockStore{
		userId: 12,
	}

	r := gin.New()
	r.Use(middleware.TestSuite())
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.POST("/games/2/asteroids/completed", CompletedAsteroid(store))

	// happy path
	body := strings.NewReader(`{"ship_id": 5}`)
	req, _ := http.NewRequest("POST", "/games/2/asteroids/completed", body)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", "12")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	c.Check(store.updated, check.Equals, true)
	c.Check(store.asteroidId, check.Equals, int64(4))
}
