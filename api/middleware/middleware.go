package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/cbarraford/cryptocades-backend/store/session"
)

// Allow api request to masquerade as any user. This is intended only for
// development purposes
func Masquerade() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Request.Header.Get("Masquerade")

		if userId != "" {
			c.Set("userId", userId)
			c.Set("escalated", true)
		}
	}
}

// Check that the request has been authorized
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get("userId")
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
			return
		}
	}
}

// Check that the request has been authorized
func EscalatedAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get("userId")
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
			return
		}
		escalated, ok := c.Get("escalated")
		if !ok || escalated == false {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Escalated privileges expired"))
			return
		}

	}
}

func Authenticate(store session.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header["Session"] != nil {
			token := c.Request.Header["Session"][0]
			id, escalated, err := store.Authenticate(token)
			if err != nil {
				log.Printf("Unable to authorize given token: %+v", token)
				return
			}
			c.Set("userId", strconv.FormatInt(id, 10))
			c.Set("escalated", escalated)
		}
	}
}
