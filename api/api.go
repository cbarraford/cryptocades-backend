package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/CBarraford/lotto/api/middleware"
	"github.com/CBarraford/lotto/api/users"
	"github.com/CBarraford/lotto/store"
)

func GetAPIService(store store.Store) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", ping())

	r.POST("/login", users.Login(store.Users, store.Sessions))
	//r.DELETE("/logout", users.Logout(store.Sessions))

	usersGroup := r.Group("/users", middleware.AuthRequired())
	{
		usersGroup.POST("/", users.Create(store.Users))
		usersGroup.GET("/:id", users.Get(store.Users))
	}
	_ = usersGroup

	return r
}

// health-check to test service is up
func ping() func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
