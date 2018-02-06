package jackpots

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/store/jackpot"
)

type JackpotListSuite struct{}

var _ = check.Suite(&JackpotListSuite{})

type mockListStore struct {
	jackpot.Dummy
}

func (*mockListStore) List() ([]jackpot.Record, error) {
	return []jackpot.Record{
		{
			Jackpot: 88,
		},
	}, nil
}

func (s *JackpotListSuite) TestList(c *check.C) {
	store := &mockListStore{}

	r := gin.New()
	r.GET("/jackpots", List(store))

	// happy path
	req, _ := http.NewRequest("GET", "/jackpots", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var records []jackpot.Record
	c.Assert(json.Unmarshal(w.Body.Bytes(), &records), check.IsNil)
	c.Assert(records, check.HasLen, 1)
	record := records[0]
	c.Check(record.Jackpot, check.Equals, float64(88))
}
