package boosts

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"

	"github.com/cbarraford/cryptocades-backend/api/middleware"
	"github.com/cbarraford/cryptocades-backend/store/boost"
)

type BoostSuite struct{}

var _ = Suite(&BoostSuite{})

type mockBoostStore struct {
	boost.Dummy
	boostId  int64
	incomeId int64
	created  bool
}

func (m *mockBoostStore) Assign(boostId, incomeId int64) error {
	m.boostId = boostId
	m.incomeId = incomeId
	m.created = true
	return nil
}

func (s *BoostSuite) TestAssign(c *C) {
	gin.SetMode(gin.ReleaseMode)

	// happy path
	store := &mockBoostStore{}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())
	r.Use(middleware.HandleErrors())
	r.PUT("/me/boosts", Assign(store))
	input := fmt.Sprintf(`{"boost_id":4, "income_id":6}`)
	body := strings.NewReader(input)
	req, _ := http.NewRequest("PUT", "/me/boosts", body)
	req.Header.Set("Masquerade", "5")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)
	c.Check(store.created, Equals, true)
	c.Check(store.boostId, Equals, int64(4))
	c.Check(store.incomeId, Equals, int64(6))
}
