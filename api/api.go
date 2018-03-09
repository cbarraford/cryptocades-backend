package api

import (
	"net/http"
	"time"

	"github.com/cbarraford/cache"
	"github.com/cbarraford/cache/persistence"
	recaptcha "github.com/ezzarghili/recaptcha-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	coinApi "github.com/miguelmota/go-coinmarketcap"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"

	"github.com/cbarraford/cryptocades-backend/api/context"
	"github.com/cbarraford/cryptocades-backend/api/facebook"
	"github.com/cbarraford/cryptocades-backend/api/games"
	"github.com/cbarraford/cryptocades-backend/api/jackpots"
	"github.com/cbarraford/cryptocades-backend/api/middleware"
	"github.com/cbarraford/cryptocades-backend/api/users"
	"github.com/cbarraford/cryptocades-backend/store"
)

func GetAPIService(store store.Store, agent newrelic.Application, captcha recaptcha.ReCAPTCHA) *gin.Engine {
	mem := persistence.NewInMemoryStore(60 * time.Second)

	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Authenticate(store.Sessions))
	r.Use(nrgin.Middleware(agent))

	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://staging-app.cryptocades.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Session"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.HandleErrors())

	r.Static("/docs", "./docs")

	r.GET("/ping", ping())
	r.GET("/currency/price/:symbol", cache.CachePage(mem, time.Minute, price()))

	r.GET("/me", middleware.AuthRequired(), users.Me(store.Users))
	r.GET("/me/balance", middleware.AuthRequired(), users.Balance(store.Incomes, store.Entries))
	r.GET("/me/incomes", middleware.AuthRequired(), users.Incomes(store.Incomes))
	r.GET("/me/incomes/rank", middleware.AuthRequired(), users.IncomeRank(store.Incomes))
	r.GET("/me/entries", middleware.AuthRequired(), users.Entries(store.Entries))
	r.PUT("/me", middleware.EscalatedAuthRequired(), users.Update(store.Users))
	// update specifically email
	r.PUT("/me/email", middleware.EscalatedAuthRequired(), users.UpdateEmail(store.Users, store.Confirmations))
	r.DELETE("/me", middleware.AuthRequired(), users.Delete(store.Users))
	r.POST("/login", users.Login(store.Users, store.Sessions))
	r.POST("/login/facebook", facebook.Login(store.Users, store.Incomes, store.Sessions))
	r.DELETE("/logout", users.Logout(store.Sessions))
	r.POST("/users", users.Create(store.Users, store.Incomes, store.Confirmations, captcha))
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

func price() func(*gin.Context) {
	return func(c *gin.Context) {
		sym := context.GetString("price", c)
		coinInfo, err := coinApi.GetCoinData(sym)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(200, gin.H{
			"usd": coinInfo.PriceUsd,
		})
	}
}
