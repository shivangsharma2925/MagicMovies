package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/controllers"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/database"
	dblogger "github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/logger"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/middleware"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/services"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/websocket"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/workers"
)

func SetupRoutes(
	router *gin.Engine,
	db *database.MongoDB,
	dbLogger *dblogger.DBLogger,
	redisOpt asynq.RedisClientOpt,
	redisClient *redis.Client,
) {
	// Initialize services with dependencies
	OtpService := services.NewOTPService(redisClient)
	EmailService := services.NewEmailService()
	movieServices := services.NewMovieService(db, dbLogger)
	jobServices := services.NewJobService(db)

	// Initialize controllers with dependencies
	movieController := controllers.NewMovieController(db, dbLogger)
	userController := controllers.NewUserController(db, dbLogger, OtpService, EmailService)
	jobController := controllers.NewJobController(db, dbLogger)

	// Worker Server
	go workers.StartWorkerServer(movieServices, jobServices, redisOpt)

	// Initialize middleware with dependencies
	authMiddleware := middleware.AuthMiddleware(dbLogger)
	ipLimiter := middleware.NewRateLimiter(2, 5)    // public APIs, rate = 2 tokens/sec, burst = 5 req at a same time
	userLimiter := middleware.NewRateLimiter(5, 10) // private APIs, rate = 5 tokens/sec, burst = 10 req at a same time

	// Special request for Websocket connection
	router.GET("/ws", websocket.HandleConnections)

	// API versioning
	api := router.Group("/api/v1")

	// Public routes
	api.POST("/register", ipLimiter.IPMiddleware(), userController.RegisterUser)
	api.POST("/login", ipLimiter.IPMiddleware(), userController.LoginUser)
	api.POST("/logout", ipLimiter.IPMiddleware(), userController.LogoutHandler)
	api.GET("/refresh", ipLimiter.IPMiddleware(), userController.RefreshTokenHandler)
	api.POST("/verify-email", ipLimiter.IPMiddleware(), userController.VerifyEmail)
	api.POST("/resend-verification", ipLimiter.IPMiddleware(), userController.ResendVerification)
	api.POST("/forgot-password", ipLimiter.IPMiddleware(), userController.ForgotPassword)
	api.POST("/reset-password", ipLimiter.IPMiddleware(), userController.ResetPassword)

	api.GET("/movies", ipLimiter.IPMiddleware(), movieController.GetMovies)
	api.GET("/genres", ipLimiter.IPMiddleware(), movieController.GetGenres)
	api.GET("/movies/suggestions", movieController.GetSuggestions)

	// Protected routes
	protected := api.Group("/")
	protected.Use(authMiddleware)
	protected.Use(userLimiter.UserMiddleware())
	{
		protected.GET("/profile/me", userController.GetProfile)

		protected.GET("/movie/:imdb_id", movieController.GetMovie)
		protected.POST("/addmovie", movieController.AddMovie)
		protected.GET("/recommendedmovies", movieController.GetRecommendedMovies)
		protected.POST("/updatereview/:imdb_id", movieController.AdminReviewUpdate)

		protected.GET("/jobs", jobController.GetJobs)
		protected.POST("/jobs/retry/:imdb_id", jobController.RetryJob)
	}
}
