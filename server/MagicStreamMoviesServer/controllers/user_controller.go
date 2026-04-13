package controllers

import (
	"context"
	"net/http"

	// "os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/database"
	dblogger "github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/logger"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/models"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// var adminPassword = os.Getenv("ADMIN_PASSWORD")

type UserController struct {
	db       *database.MongoDB
	dbLogger *dblogger.DBLogger
}

func NewUserController(db *database.MongoDB, dbLogger *dblogger.DBLogger) *UserController {
	return &UserController{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (uc *UserController) RegisterUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validations failed", "details": err.Error()})
		return
	}

	hashedPassword, err := utilities.HashPassword(user.Password)
	if err != nil {
		uc.dbLogger.Alerts("ERROR", "Password hashing failed", gin.H{
			"endpoint": "/RegisterUser",
			"user_id":  user.UserID,
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to Hash password"})
		return
	}

	usercollection := uc.db.Collection("users")

	count, err := usercollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		uc.dbLogger.Alerts("ERROR", "User existence check failed", gin.H{
			"endpoint": "/RegisterUser",
			"user_id":  user.UserID,
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	user.UserID = primitive.NewObjectID().Hex()
	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = usercollection.InsertOne(ctx, user)
	if err != nil {
		uc.dbLogger.Alerts("ERROR", "User insert failed", gin.H{
			"endpoint": "/RegisterUser",
			"user_id":  user.UserID,
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Account creation failed"})
		return
	}

	uc.dbLogger.Log("ACCOUNT", "User registered successfully", gin.H{
		"endpoint": "/RegisterUser",
		"user_id":  user.UserID,
	})

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

func (uc *UserController) LoginUser(c *gin.Context) {
	var userLoginData models.UserLogin

	if err := c.ShouldBindJSON(&userLoginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	if err := validate.Struct(userLoginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to Validate", "details": err.Error()})
		return
	}

	usercollection := uc.db.Collection("users")

	var user models.User

	err := usercollection.FindOne(ctx, bson.M{"email": userLoginData.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err = utilities.ComparePassword(user.Password, userLoginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, refreshToken, err := utilities.GenerateAllTokens(user.Email, user.FirstName, user.LastName, user.Role, user.UserID)
	if err != nil {
		uc.dbLogger.Alerts("ERROR", "Token generation failed", gin.H{
			"endpoint": "/LoginUser",
			"user_id":  user.UserID,
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	err = uc.UpdateAllTokens(user.UserID, refreshToken, c)
	if err != nil {
		uc.dbLogger.Alerts("ERROR", "Token update failed", gin.H{
			"endpoint": "/LoginUser",
			"user_id":  user.UserID,
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	uc.dbLogger.Log("ACCOUNT", "User logged in successfully", gin.H{
		"endpoint": "/LoginUser",
		"user_id":  user.UserID,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name: "access_token",
		Value: token,
		Path: "/",
		MaxAge: 1800,
		Secure: true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name: "refresh_token",
		Value: refreshToken,
		Path: "/",
		MaxAge: 24*60*60,
		Secure: true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	c.JSON(http.StatusOK, models.UserResponse{
		Email:           user.Email,
		UserId:          user.UserID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Role:            user.Role,
		FavouriteGenres: user.FavouriteGenres,
	})
}

func (uc *UserController) LogoutHandler(c *gin.Context) {
	// Clear the access_token cookie

	var UserLogout struct {
		UserId string `json:"user_id"`
	}

	err := c.ShouldBindJSON(&UserLogout)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err = uc.UpdateAllTokens(UserLogout.UserId, "", c) // Clear tokens in the database

	// Optionally, you can also remove the user session from the database if needed

	if err != nil {
		uc.dbLogger.Alerts("ERROR", "Error logging out", gin.H{
			"endpoint": "/LogoutHandler",
			"user_id":  UserLogout.UserId,
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging out"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:  "access_token",
		Value: "",
		Path:  "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	uc.dbLogger.Log("ACCOUNT", "User logged out successfully", gin.H{
		"endpoint": "/LogoutHandler",
		"user_id":  UserLogout.UserId,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (uc *UserController) RefreshTokenHandler(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(c, 10*time.Second)
	defer cancel()

	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve refresh token from cookie"})
		return
	}

	claim, err := utilities.ValidateRefreshToken(refreshToken)
	if err != nil || claim == nil {
		uc.dbLogger.Alerts("ERROR", "refresh token validation failed", gin.H{
			"endpoint": "/RefreshTokenHandler",
			"error":    err.Error(),
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	usercollection := uc.db.Collection("users")

	var user models.User
	err = usercollection.FindOne(ctx, bson.D{{Key: "user_id", Value: claim.UserId}}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	newToken, newRefreshToken, _ := utilities.GenerateAllTokens(user.Email, user.FirstName, user.LastName, user.Role, user.UserID)
	err = uc.UpdateAllTokens(user.UserID, newRefreshToken, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating tokens"})
		return
	}

	c.SetCookie("access_token", newToken, 1800, "/", "localhost", true, true)          // expires in 30 mins
	c.SetCookie("refresh_token", newRefreshToken, 86400, "/", "localhost", true, true) // expires in 24 hours

	c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed"})
}

func (uc *UserController) UpdateAllTokens(userid, refreshtoken string, c *gin.Context) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	update_at := time.Now()

	updates := bson.M{
		"$set": bson.M{
			"refresh_token": refreshtoken,
			"updated_at":    update_at,
		},
	}

	usercollection := uc.db.Collection("users")

	_, err := usercollection.UpdateOne(ctx, bson.M{"user_id": userid}, updates)

	if err != nil {
		return err
	}

	return nil
}
