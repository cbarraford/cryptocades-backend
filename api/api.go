package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/CBarraford/lotto/api/jackpots"
	"github.com/CBarraford/lotto/api/middleware"
	"github.com/CBarraford/lotto/api/users"
	"github.com/CBarraford/lotto/store"
)

func GetAPIService(store store.Store) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Authenticate(store.Sessions))

	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Session"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/docs", "./docs")

	r.GET("/ping", ping())
	r.GET("/me", users.Me(store.Users))
	r.PUT("/me", users.Update(store.Users))
	r.POST("/login", users.Login(store.Users, store.Sessions))
	r.DELETE("/logout", users.Logout(store.Sessions))
	r.POST("/users", users.Create(store.Users))

	usersGroup := r.Group("/users", middleware.AuthRequired())
	{
		usersGroup.GET("/:id", users.Get(store.Users))
	}

	jackpotsGroup := r.Group("/jackpots")
	{
		jackpotsGroup.GET("/", jackpots.List(store.Jackpots))
		jackpotsGroup.GET("/:id/odds", jackpots.Odds(store.Entries))
		jackpotsGroup.POST("/:id/enter", jackpots.Enter(store.Users, store.Entries))
	}

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
