package context

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"

	"github.com/CBarraford/lotto/api/middleware"
)

func TestPackage(t *testing.T) { TestingT(t) }

type ContextSuite struct{}

var _ = Suite(&ContextSuite{})

func (s *ContextSuite) TestContext(c *C) {
	r := gin.New()
	r.Use(middleware.Masquerade())
	r.Use(middleware.AuthRequired())

	masqueradeId := "5"
	r.GET("/test/:id", func(context *gin.Context) {
		requestId, err := GetUserId(context)
		c.Assert(err, IsNil)
		c.Check(requestId, Equals, int64(5), Commentf("RequestId: %s", requestId))

		targetId, err := GetInt64("id", context)
		c.Assert(err, IsNil)
		c.Check(targetId, Equals, int64(8))
		context.String(200, "OK")
	})

	// happy path
	req, _ := http.NewRequest("GET", "/test/8", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Masquerade", masqueradeId)
	r.ServeHTTP(w, req)
	c.Assert(w.Code, Equals, 200)
}
