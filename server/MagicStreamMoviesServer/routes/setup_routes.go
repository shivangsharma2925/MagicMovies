package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/controllers"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/database"
	dblogger "github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/logger"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/middleware"
)

func SetupRoutes(
	router *gin.Engine,
	db *database.MongoDB,
	dbLogger *dblogger.DBLogger,
) {

	// Initialize controllers with dependencies
	movieController := controllers.NewMovieController(db, dbLogger)
	userController := controllers.NewUserController(db, dbLogger)
	
	// Initialize middleware with dependencies
	authMiddleware := middleware.AuthMiddleware(dbLogger)
	ipLimiter := middleware.NewRateLimiter(2, 5)    // public APIs, rate = 2 tokens/sec, burst = 5 req at a same time
	userLimiter := middleware.NewRateLimiter(5, 10) // private APIs, rate = 5 tokens/sec, burst = 10 req at a same time

	// API versioning
	api := router.Group("/api/v1")

	// Public routes
	api.POST("/register", ipLimiter.IPMiddleware(), userController.RegisterUser)
	api.POST("/login", ipLimiter.IPMiddleware(), userController.LoginUser)
	api.POST("/logout", ipLimiter.IPMiddleware(), userController.LogoutHandler)
	api.GET("/refresh", ipLimiter.IPMiddleware(), userController.RefreshTokenHandler)
	
	api.GET("/movies", ipLimiter.IPMiddleware(), movieController.GetMovies)
	api.GET("/genres", ipLimiter.IPMiddleware(), movieController.GetGenres)

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
	}
}
