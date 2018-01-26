package games

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store/game"
)

func TestPackage(t *testing.T) { check.TestingT(t) }

type GameListSuite struct{}

var _ = check.Suite(&GameListSuite{})

func (s *GameListSuite) TestList(c *check.C) {
	store := game.NewStore()

	r := gin.New()
	r.GET("/games", List(store))

	// happy path
	req, _ := http.NewRequest("GET", "/games", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var games []game.Record
	c.Assert(json.Unmarshal(w.Body.Bytes(), &games), check.IsNil)
	c.Assert(games, check.HasLen, 2)
	c.Check(games[0].Id, check.Equals, 1)
	c.Check(games[0].Name, check.Equals, "Tallest Tower")
	c.Check(games[1].Id, check.Equals, 2)
	c.Check(games[1].Name, check.Equals, "Asteroid Tycoon")
}
