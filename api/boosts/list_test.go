package boosts

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/api/middleware"
	"github.com/cbarraford/cryptocades-backend/store/boost"
)

type UserBoostsSuite struct{}

var _ = check.Suite(&UserBoostsSuite{})

type mockBoostsBoostsStore struct {
	boost.Dummy
}

func (*mockBoostsBoostsStore) ListByUser(id int64) ([]boost.Record, error) {
	return []boost.Record{
		{Id: 15, IncomeId: 4, UserId: id},
	}, nil
}

func (s *UserBoostsSuite) TestBoosts(c *check.C) {
	store := &mockBoostsBoostsStore{}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())

	r.GET("/me/boosts", ListByUser(store))

	// happy path
	req, _ := http.NewRequest("GET", "/me/boosts", nil)
	req.Header.Set("Masquerade", "12")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var records []boost.Record
	c.Assert(json.Unmarshal(w.Body.Bytes(), &records), check.IsNil)
	c.Assert(records, check.HasLen, 1)
	c.Check(records[0].Id, check.Equals, int64(15))
	c.Check(records[0].IncomeId, check.Equals, int64(4))
	c.Check(records[0].UserId, check.Equals, int64(12))
}
