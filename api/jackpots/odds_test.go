package jackpots

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"

	"github.com/CBarraford/lotto/api/middleware"
	"github.com/CBarraford/lotto/store/entry"
)

type JackpotOddsSuite struct{}

var _ = check.Suite(&JackpotOddsSuite{})

type mockOddsStore struct {
	entry.Dummy
}

func (*mockOddsStore) GetOdds(j, i int64) (entry.Odds, error) {
	return entry.Odds{
		JackpotId: j,
		Total:     45,
		Entries:   12,
	}, nil
}

func (s *JackpotOddsSuite) TestOdds(c *check.C) {
	store := &mockOddsStore{}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.GET("/jackpots/:id/odds", Odds(store))

	// happy path
	req, _ := http.NewRequest("GET", "/jackpots/23/odds", nil)
	req.Header.Set("Masquerade", "44")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var odd entry.Odds
	c.Assert(json.Unmarshal(w.Body.Bytes(), &odd), check.IsNil)
	c.Check(odd.JackpotId, check.Equals, int64(23))
	c.Check(odd.Total, check.Equals, int64(45))
	c.Check(odd.Entries, check.Equals, int64(12))
}
