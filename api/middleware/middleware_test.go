package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type MiddlewareSuite struct{}

var _ = Suite(&MiddlewareSuite{})

func (s *MiddlewareSuite) TestMasquerade(c *C) {
	r := gin.New()
	r.Use(Masquerade())
	r.Use(AuthRequired())

	masqueradeId := "myId"
	r.GET("/test", func(context *gin.Context) {
		playerId, _ := context.Get("playerId")
		c.Check(playerId, Equals, masqueradeId, Commentf("PlayerId: %s", playerId))
		context.String(200, "OK")
	})

	// happy path
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", masqueradeId)
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)

	// no auth in request
	req, _ = http.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 401)
}
