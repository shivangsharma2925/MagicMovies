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
	movieServices := services.NewMovieService(db, dbLogger)
	jobServices := services.NewJobService(db)

	// Initialize controllers with dependencies
	movieController := controllers.NewMovieController(db, dbLogger)
	userController := controllers.NewUserController(db, dbLogger, OtpService)
	jobController := controllers.NewJobController(db, dbLogger)
	watchListController := controllers.NewWatchListController(db, dbLogger)

	// Worker Server
	go workers.StartWorkerServer(movieServices, jobServices, redisOpt)

	// Initialize middleware with dependencies
	authMiddleware := middleware.AuthMiddleware(dbLogger)
	ipLimiter := middleware.NewRateLimiter(2, 5)    // public APIs, rate = 2 tokens/sec, burst = 5 req at a same time
	userLimiter := middleware.NewRateLimiter(5, 10) // private APIs, rate = 5 tokens/sec, burst = 10 req at a same time

	// router.OPTIONS("/*path", func(c *gin.Context) {
	// 	c.AbortWithStatus(204)
	// })

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
	api.GET("/movies/suggestions", ipLimiter.IPMiddleware(), movieController.GetSuggestions)
	api.GET("/movies/heroMovies", ipLimiter.IPMiddleware(), movieController.GetHeroMovies)
	api.GET("/movie/:imdb_id", ipLimiter.IPMiddleware(), movieController.GetMovie)
	api.GET("/category/:title", ipLimiter.IPMiddleware(), movieController.GetByCategory)

	// Used when some bulk updates on movie collection needs to be done
	// api.GET("/UpdateMovies", ipLimiter.IPMiddleware(), movieServices.UpdateMovies)

	// Protected routes
	protected := api.Group("/")
	protected.Use(authMiddleware)
	protected.Use(userLimiter.UserMiddleware())
	{
		protected.GET("/profile/me", userController.GetProfile)
		protected.POST("/delete-account", userController.DeleteAccount)

		protected.POST("/addmovie", movieController.AddMovie)
		protected.POST("/addmovieai", movieController.AddMovieViaAI)
		protected.GET("/recommendedmovies", movieController.GetRecommendedMovies)
		protected.POST("/updatereview/:imdb_id", movieController.AdminReviewUpdate)

		protected.GET("/jobs", jobController.GetJobs)
		protected.POST("/jobs/retry/:imdb_id", jobController.RetryJob)
		protected.POST("/ai/imdb", jobController.AskAiImdb)

		protected.GET("/watchlists", watchListController.GetWatchLists)
		protected.POST("/watchlists", watchListController.NewWatchList)
		protected.DELETE("/watchlists/:watchlist_id", watchListController.DeleteWatchList)
		protected.POST("/watchlists/:watchlist_id", watchListController.RenameWatchLists)
		protected.POST("/watchlists/:watchlist_id/movies", watchListController.AddMovietoWatchList)
		protected.DELETE("/watchlists/:watchlist_id/movies/:movie_id", watchListController.DeleteMoviefromWatchList)
		protected.POST("/watchlists/:watchlist_id/movies/:movie_id/watched", watchListController.ToggleWatch)

	}
}
