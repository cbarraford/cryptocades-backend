package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"

	"github.com/CBarraford/lotto/api/middleware"
	"github.com/CBarraford/lotto/store/user"
)

type UserMeSuite struct{}

var _ = check.Suite(&UserMeSuite{})

type mockMeStore struct {
	user.Dummy
}

func (*mockMeStore) Get(id int64) (user.Record, error) {
	return user.Record{
		Id:       id,
		Username: "bob",
	}, nil
}

func (s *UserMeSuite) TestMe(c *check.C) {
	store := &mockMeStore{}

	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())

	r.GET("/users/me", Me(store))

	// happy path
	req, _ := http.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Masquerade", "12")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	record := user.Record{}
	c.Assert(json.Unmarshal(w.Body.Bytes(), &record), check.IsNil)
	c.Check(record.Id, check.Equals, int64(12))
	c.Check(record.Username, check.Equals, "bob")
}
