package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CBarraford/lotto/store/session"
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

func Authenticate(store session.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header["Session"] != nil {
			token := c.Request.Header["Session"][0]
			id, err := store.Authenticate(token)
			if err != nil {
				log.Printf("Unable to authorize given token: %+v", token)
				return
			}
			c.Set("userId", id)
		}
	}
}
