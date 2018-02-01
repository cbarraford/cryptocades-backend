package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/api/games"
	"github.com/cbarraford/cryptocades-backend/api/jackpots"
	"github.com/cbarraford/cryptocades-backend/api/middleware"
	"github.com/cbarraford/cryptocades-backend/api/users"
	"github.com/cbarraford/cryptocades-backend/store"
)

func GetAPIService(store store.Store, agent newrelic.Application) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Authenticate(store.Sessions))
	r.Use(middleware.HandleErrors())
	r.Use(nrgin.Middleware(agent))

	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://lotto-ui-staging.herokuapp.com", "https://staging.cryptokade.com", "https://staging.cryptocades.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Session"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/docs", "./docs")

	r.GET("/ping", ping())
	r.GET("/me", middleware.AuthRequired(), users.Me(store.Users))
	r.GET("/me/balance", middleware.AuthRequired(), users.Balance(store.Incomes, store.Entries))
	r.GET("/me/incomes", middleware.AuthRequired(), users.Incomes(store.Incomes))
	r.GET("/me/entries", middleware.AuthRequired(), users.Entries(store.Entries))
	r.PUT("/me", middleware.EscalatedAuthRequired(), users.Update(store.Users))
	// update specifically email
	r.PUT("/me/email", middleware.EscalatedAuthRequired(), users.UpdateEmail(store.Users, store.Confirmations))
	r.DELETE("/me", middleware.AuthRequired(), users.Delete(store.Users))
	r.POST("/login", users.Login(store.Users, store.Sessions))
	r.DELETE("/logout", users.Logout(store.Sessions))
	r.POST("/users", users.Create(store.Users, store.Confirmations))
	r.POST("/users/confirmation/:code",
		users.Confirm(store.Confirmations, store.Users),
	)
	r.POST("/users/password_reset",
		users.PasswordResetInit(store.Confirmations, store.Users),
	)
	r.POST("/users/password_reset/:code",
		users.PasswordReset(store.Confirmations, store.Users),
	)

	usersGroup := r.Group("/users", middleware.AuthRequired())
	{
		usersGroup.GET("/:id", users.Get(store.Users))
	}

	r.GET("/games", games.List(store.Games))

	jackpotsGroup := r.Group("/jackpots")
	{
		jackpotsGroup.GET("/", jackpots.List(store.Jackpots))
		jackpotsGroup.GET("/:id/odds", jackpots.Odds(store.Entries))
		jackpotsGroup.POST("/:id/enter", jackpots.Enter(store.Entries, store.Users, store.Jackpots))
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
