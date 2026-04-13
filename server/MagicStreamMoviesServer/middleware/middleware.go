package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dblogger "github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/logger"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/utilities"
)

func AuthMiddleware(logger *dblogger.DBLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utilities.GetAccessToken(c)

		if err != nil {
			logger.Alerts("ERROR", "Missing or invalid auth header", gin.H{
				"endpoint": "AuthMiddleware",
				"error":    err.Error(),
			})
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		claims, err := utilities.ValidateToken(token)
		if err != nil {
			logger.Alerts("ERROR", "Invalid token", gin.H{
				"endpoint": "AuthMiddleware",
				"error":    err.Error(),
			})
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userId", claims.UserId)
		c.Set("role", claims.Role)

		go func ()  {
			logger.Log("ACCOUNT", "Request validated", gin.H{
				"endpoint": "AuthMiddleware",
				"Ip":       c.ClientIP(),
				"UseriD":   claims.UserId,
			})
		}()

		c.Next()
	}
}
