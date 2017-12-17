package context

import (
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
