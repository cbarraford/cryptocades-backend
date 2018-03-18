package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stvp/rollbar"

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

func AdminAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Get("token")
		// TODO: hard coded token key for admin api access (ie for bots to
		// authenticate)
		if token == "QieDpVTtcnBgFVDPccRmDa98" {
			return
		}
		_, ok := c.Get("userId")
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
			return
		}
		admin, ok := c.Get("admin")
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Access Denied"))
			return
		}
		if !admin.(bool) {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Access Denied"))
			return
		}
	}
}

func Authenticate(store session.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header["Session"] != nil {
			token := c.Request.Header["Session"][0]
			id, escalated, admin, err := store.Authenticate(token)
			if err != nil {
				log.Printf("Unable to authorize given token: %+v", token)
				return
			}
			c.Set("userId", strconv.FormatInt(id, 10))
			c.Set("admin", admin)
			c.Set("escalated", escalated)
		}
		if c.Request.Header["Token"] != nil {
			token := c.Request.Header["Token"][0]
			c.Set("token", token)
		}
	}
}

func HandleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // execute all the handlers

		errorToPrint := c.Errors.Last()
		if errorToPrint != nil {
			rollbar.Error(rollbar.ERR, errorToPrint)
			c.JSON(-1, gin.H{
				"message": errorToPrint.Error(),
			})
		}
	}
}
