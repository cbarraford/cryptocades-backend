package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var httpClient = "httpClient"

// Allow api request to masquerade as any user. This is intended only for
// development purposes
func Masquerade() gin.HandlerFunc {
	return func(c *gin.Context) {
		playerId := c.Request.Header.Get("Masquerade")

		if playerId != "" {
			c.Set("playerId", playerId)
		}
	}
}

// Check that the request has been authorized
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get("playerId")
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}
	}
}
