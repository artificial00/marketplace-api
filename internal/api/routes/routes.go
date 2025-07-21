package routes

import (
	"database/sql"
	"marketplace-api/internal/database"

	"github.com/gin-gonic/gin"
	"marketplace-api/internal/api/handlers"
	"marketplace-api/internal/config"
	"marketplace-api/internal/database/postgres"
	"marketplace-api/internal/service"
	"marketplace-api/pkg/middleware"
)

// SetupRoutes настраивает все маршруты приложения согласно ТЗ
func SetupRoutes(router *gin.Engine, db *sql.DB, cfg *config.Config) {
	router.Use(middleware.CORSMiddleware())

	userRepo := postgres.NewUserRepository(db)
	listingRepo := postgres.NewListingRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWT.Secret)
	listingService := service.NewListingService(listingRepo)

	authHandler := handlers.NewAuthHandler(authService)
	listingHandler := handlers.NewListingHandler(listingService)

	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			if err := database.CheckConnection(db); err != nil {
				c.JSON(500, gin.H{
					"status":   "unhealthy",
					"database": "disconnected",
					"error":    err.Error(),
				})
				return
			}

			c.JSON(200, gin.H{
				"status":   "healthy",
				"database": "connected",
				"service":  "marketplace-api",
				"version":  "1.0.0",
			})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		listings := api.Group("/listings")
		{
			listings.GET("/", listingHandler.GetListings)
			listings.GET("/:id", listingHandler.GetListing)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			protected.GET("/auth/me", authHandler.Me)

			protectedListings := protected.Group("/listings")
			{
				protectedListings.POST("/", listingHandler.CreateListing)
				protectedListings.GET("/my", listingHandler.GetMyListings)
				protectedListings.PUT("/:id", listingHandler.UpdateListing)
				protectedListings.DELETE("/:id", listingHandler.DeleteListing)
			}
		}
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Marketplace API",
			"version": "1.0.0",
			"status":  "running",
			"endpoints": gin.H{
				"health":   "/api/health",
				"auth":     "/api/auth/*",
				"listings": "/api/listings/*",
			},
		})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "not_found",
			"message": "The requested resource was not found",
			"path":    c.Request.URL.Path,
		})
	})
}
