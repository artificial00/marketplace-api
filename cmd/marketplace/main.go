package main

import (
	"context"
	"fmt"
	"marketplace-api/internal/api/routes"
	"marketplace-api/internal/config"
	"marketplace-api/internal/database"
	"marketplace-api/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// @title Marketplace API
// @version 1.0
// @description REST API для маркетплейса с авторизацией и работой с объявлениями
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme. Example: "Authorization: Bearer {token}"

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	log := logger.New()
	log.Info("Starting marketplace API server")
	log.Info("Configuration loaded successfully")

	db, err := database.NewConnection(cfg.GetDatabaseURL())
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	log.Info("Database connection established")

	if err := database.CreateTables(db); err != nil {
		log.Error("Failed to create database tables", "error", err)
		os.Exit(1)
	}

	log.Info("Database tables initialized")

	gin.SetMode(gin.DebugMode)
	router := gin.New()

	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: os.Stdout,
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("[%s] %s %s %d %s %s\n",
				param.TimeStamp.Format("2006-01-02 15:04:05"),
				param.Method,
				param.Path,
				param.StatusCode,
				param.Latency,
				param.ErrorMessage,
			)
		},
	}))

	router.Use(gin.Recovery())

	routes.SetupRoutes(router, db, cfg)

	server := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("Server starting",
			"address", cfg.GetServerAddress(),
			"health_check", fmt.Sprintf("http://%s/api/health", cfg.GetServerAddress()),
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	log.Info("Marketplace API is running")
	log.Info("API Documentation will be available at: http://" + cfg.GetServerAddress() + "/swagger/index.html")
	log.Info("Health Check: http://" + cfg.GetServerAddress() + "/api/health")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Server stopped gracefully")
}
