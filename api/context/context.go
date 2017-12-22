package context

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetInt64(param string, c *gin.Context) (id int64, err error) {
	id, err = strconv.ParseInt(c.Param(param), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("URL id must be a number"))
	}
	return
}

func GetUserId(c *gin.Context) (id int64, err error) {
	userId, ok := c.Get("userId")
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, errors.New("Unauthorized"))
		return 0, errors.New("Unauthorized")
	}
	id, err = strconv.ParseInt(userId.(string), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("Bad authenticating user id"))
		return 0, errors.New("Bad authenticating user id")
	}
	return
}
