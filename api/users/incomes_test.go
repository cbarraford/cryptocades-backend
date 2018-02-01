package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/api/middleware"
	"github.com/cbarraford/cryptocades-backend/store/income"
)

type UserIncomesSuite struct{}

var _ = check.Suite(&UserIncomesSuite{})

type mockIncomesStore struct {
	income.Dummy
}

func (*mockIncomesStore) ListByUser(id int64) ([]income.Record, error) {
	return []income.Record{
		{Id: 15, GameId: 4, UserId: id, Amount: 45},
	}, nil
}

func (s *UserIncomesSuite) TestIncomes(c *check.C) {
	store := &mockIncomesStore{}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())

	r.GET("/me/incomes", Incomes(store))

	// happy path
	req, _ := http.NewRequest("GET", "/me/incomes", nil)
	req.Header.Set("Masquerade", "12")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	var records []income.Record
	c.Assert(json.Unmarshal(w.Body.Bytes(), &records), check.IsNil)
	c.Assert(records, check.HasLen, 1)
	c.Check(records[0].Id, check.Equals, int64(15))
	c.Check(records[0].GameId, check.Equals, int64(4))
	c.Check(records[0].UserId, check.Equals, int64(12))
	c.Check(records[0].Amount, check.Equals, 45)
}
