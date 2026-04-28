package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	// "github.com/gin-contrib/cors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/database"
	dblogger "github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/logger"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/routes"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/services"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/workers"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load config (env)
	port := getEnv("PORT", "8080")

	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Set Gin mode (important for production), this removed logging and extra details thus taking less resources
	gin.SetMode(gin.ReleaseMode)

	// creates a new router instance with two middleware already attached, logging and recovery. This will create middlewares with default behaviour, so New() gived more manual access
	// router := gin.Default()
	router := gin.New()

	// Custom middleware
	router.Use(gin.Recovery())

	// Load env
	err := godotenv.Load(".env")
	if err != nil {
		logger.Warn("No .env file found")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DATABASE_NAME")

	if mongoURI == "" || dbName == "" {
		logger.Error("Missing env variables")
		os.Exit(1)
	}

	db, err := database.NewMongoDB(context.Background(), mongoURI, dbName, logger)
	if err != nil {
		logger.Error("DB init failed", "error", err)
		os.Exit(1)
	}

	if err := database.InitDB(db, logger); err != nil {
		logger.Warn("DB init failed", "error", err)
	}

	var dbLogger *dblogger.DBLogger
	dbLogger = dblogger.NewDBLogger(db, logger)

	// Initialize services with dependencies
	movieServices := services.NewMovieService(db, dbLogger)

	//Initialize 3 parallel worker
	workers.StartWorkers(ctx, movieServices, 3)

	//CORS policy
	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	AllowCredentials: true,
	// }))

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	var origins []string
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
	} else {
		origins = []string{"http://localhost:5173"}
	}

	config := cors.Config{}
	config.AllowOrigins = origins
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))

	// establish routes for requests
	routes.SetupRoutes(router, db, dbLogger)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "running",
		})
	})

	// router.Run() method kills the server whenever fatal occurs thus leading to faulty or dirty data since requests abrubtly terminates
	// err := router.Run(":8080")
	// Create server manually (needed for graceful shutdown)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Run server in goroutine in order to avoid blocking the code
	go func() {
		logger.Info("Server starting", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)                      // Creates a channel to listen for OS signals.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // this  Tells the app to wait until you press Ctrl+C (SIGINT) or the cloud provider sends a stop command (SIGTERM).

	<-quit // This line blocks. The code sits here and waits until a signal is received.

	logger.Info("Shutting down server...")

	// stop workers
	cancel()

	// 5 seconds are given for all the pending requests to complete after that server will terminate them forcefully
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Forced to shutdown", "error", err)
	}

	logger.Info("Server exited properly")
}

// Helper for env variables
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
