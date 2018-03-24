package matchups

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/api/middleware"
	"github.com/cbarraford/cryptocades-backend/store/matchup"
)

func TestPackage(t *testing.T) { check.TestingT(t) }

type MatchupSuite struct{}

var _ = check.Suite(&MatchupSuite{})

type mockMatchupStore struct {
	matchup.Dummy
}

func (*mockMatchupStore) Get(e string, o int, id int64) (matchup.Record, error) {
	return matchup.Record{
		UserId:   id,
		Username: "bob",
		Rank:     5,
		Score:    12,
	}, nil
}

func (*mockMatchupStore) GetTopPerformers(s string, o int, t int) ([]matchup.Record, error) {
	return []matchup.Record{
		{UserId: 1, Rank: 1, Score: 1, Username: "Larry"},
		{UserId: 2, Rank: 2, Score: 2, Username: "Moe"},
		{UserId: 3, Rank: 3, Score: 3, Username: "Curly"},
	}, nil
}

func (s *MatchupSuite) TestGet(c *check.C) {
	store := &mockMatchupStore{}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.Use(middleware.HandleErrors())

	r.GET("/matchups/:event/:offset/me", Get(store))

	// happy path
	req, _ := http.NewRequest("GET", "/matchups/daily/0/me", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", "8")
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	record := matchup.Record{}
	c.Assert(json.Unmarshal(w.Body.Bytes(), &record), check.IsNil)
	c.Check(record.Username, check.Equals, "bob")
	c.Check(record.Rank, check.Equals, 5)
	c.Check(record.Score, check.Equals, 12)
}

func (s *MatchupSuite) TestList(c *check.C) {
	store := &mockMatchupStore{}

	r := gin.New()
	r.Use(middleware.HandleErrors())

	r.GET("/matchups/:event/:offset", List(store))

	// happy path
	req, _ := http.NewRequest("GET", "/matchups/daily/0", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	records := []matchup.Record{}
	c.Assert(json.Unmarshal(w.Body.Bytes(), &records), check.IsNil)
	c.Assert(records, check.HasLen, 3)
	c.Check(records[0].Username, check.Equals, "Larry")
	c.Check(records[0].Rank, check.Equals, 1)
	c.Check(records[0].Score, check.Equals, 1)
	c.Check(records[1].Username, check.Equals, "Moe")
	c.Check(records[1].Rank, check.Equals, 2)
	c.Check(records[1].Score, check.Equals, 2)
	c.Check(records[2].Username, check.Equals, "Curly")
	c.Check(records[2].Rank, check.Equals, 3)
	c.Check(records[2].Score, check.Equals, 3)
}
