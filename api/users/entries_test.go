package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"

	"github.com/CBarraford/lotto/api/middleware"
	"github.com/CBarraford/lotto/store/entry"
)

type UserEntriesSuite struct{}

var _ = check.Suite(&UserEntriesSuite{})

type mockEntriesEntriesStore struct {
	entry.Dummy
}

func (*mockEntriesEntriesStore) ListByUser(id int64) ([]entry.Record, error) {
	return []entry.Record{
		{Id: 15, JackpotId: 4, UserId: id, Amount: 45},
	}, nil
}

func (s *UserEntriesSuite) TestEntries(c *check.C) {
	store := &mockEntriesEntriesStore{}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())

	r.GET("/me/entries", Entries(store))

	// happy path
	req, _ := http.NewRequest("GET", "/me/entries", nil)
	req.Header.Set("Masquerade", "12")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var records []entry.Record
	c.Assert(json.Unmarshal(w.Body.Bytes(), &records), check.IsNil)
	c.Assert(records, check.HasLen, 1)
	c.Check(records[0].Id, check.Equals, int64(15))
	c.Check(records[0].JackpotId, check.Equals, int64(4))
	c.Check(records[0].UserId, check.Equals, int64(12))
	c.Check(records[0].Amount, check.Equals, 45)
}
