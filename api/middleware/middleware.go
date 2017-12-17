package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Allow api request to masquerade as any user. This is intended only for
// development purposes
func Masquerade() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Request.Header.Get("Masquerade")

		if userId != "" {
			c.Set("userId", userId)
		}
	}
}

// Check that the request has been authorized
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get("userId")
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}
	}
}
