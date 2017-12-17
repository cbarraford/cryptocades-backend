package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	check "gopkg.in/check.v1"

	"github.com/CBarraford/lotto/store/user"
)

type UserGetSuite struct{}

var _ = check.Suite(&UserGetSuite{})

type mockGetStore struct {
	user.Dummy
}

func (*mockGetStore) Get(id int64) (user.Record, error) {
	return user.Record{
		Id:       id,
		Username: "bob",
	}, nil
}

func (s *UserGetSuite) TestGet(c *check.C) {
	store := &mockGetStore{}

	r := gin.New()
	r.GET("/users/:id", Get(store))

	// happy path
	req, _ := http.NewRequest("GET", "/users/12", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, check.Equals, 200)

	record := user.Record{}
	c.Assert(json.Unmarshal(w.Body.Bytes(), &record), check.IsNil)
	c.Check(record.Id, check.Equals, int64(12))
	c.Check(record.Username, check.Equals, "bob")
}
