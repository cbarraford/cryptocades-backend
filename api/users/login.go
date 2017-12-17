package users

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CBarraford/lotto/store/session"
	"github.com/CBarraford/lotto/store/user"
)

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(store user.Store, sessionStore session.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		record := user.Record{}
		var json login
		err = c.BindJSON(&json)
		if err == nil {
			record, err = store.Authenticate(json.Username, json.Password)
			if err != nil {
				log.Printf("ErrA: %+v", err)
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
		} else {
			log.Printf("ErrA: parse json")
			c.AbortWithError(http.StatusBadRequest, errors.New("Could not parse json body"))
			return
		}

		sessionRecord := session.Record{
			UserId: record.Id,
		}
		err = sessionStore.Create(&sessionRecord, 30)
		if err != nil {
			log.Printf("ErrB: %+v", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, sessionRecord)
	}
}
