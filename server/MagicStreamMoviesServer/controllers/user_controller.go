package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	// "os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/database"
	dblogger "github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/logger"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/models"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/services"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var adminPassword = os.Getenv("ADMIN_PASSWORD")

type UserController struct {
	db         *database.MongoDB
	dbLogger   *dblogger.DBLogger
	otpService *services.OTPService
}

func NewUserController(db *database.MongoDB, dbLogger *dblogger.DBLogger, otpService *services.OTPService) *UserController {
	return &UserController{
		db:         db,
		dbLogger:   dbLogger,
		otpService: otpService,
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

	adminDowngraded := false
	if user.Role == "ADMIN" {
		if user.AdminPassword != os.Getenv("ADMIN_PASSWORD") {
			user.Role = "USER"
			adminDowngraded = true
		}
	}

	user.UserID = primitive.NewObjectID().Hex()
	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsVerified = false
	if len(user.FirstName) > 0 {
		user.FirstName = strings.ToUpper(user.FirstName[:1]) + user.FirstName[1:]
	}
	if len(user.LastName) > 0 {
		user.LastName = strings.ToUpper(user.LastName[:1]) + user.LastName[1:]
	}

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

	sendEmail := true

	otp, err := uc.otpService.GenerateSecureOTP()
	if err != nil {
		uc.dbLogger.Alerts("ERROR", "OTP generation failed", gin.H{
			"endpoint":   "/RegisterUser",
			"user_email": user.Email,
			"error":      err.Error(),
		})
		sendEmail = false
	}

	if sendEmail {
		ctxOtp, cancelOTP := context.WithTimeout(c, 20*time.Second)
		defer cancelOTP()
		err = uc.otpService.SaveVerificationOTP(ctxOtp, user.UserID, otp)

		if err != nil {
			uc.dbLogger.Alerts("ERROR", "Error in saving OTP", gin.H{
				"endpoint":   "/RegisterUser",
				"user_email": user.Email,
				"error":      err.Error(),
			})
			sendEmail = false
		}

		if sendEmail {
			err = services.SendVerificationOTP(user.Email, otp)
			if err != nil {
				uc.dbLogger.Alerts("ERROR", "Error sending email", gin.H{
					"endpoint":   "/RegisterUser",
					"user_email": user.Email,
					"error":      err.Error(),
				})
				sendEmail = false
			} else {
				uc.otpService.SetResendCooldown(ctxOtp, user.UserID)
			}
		}
	}

	verificationMessage := "Account created successfully. Redirecting to Verify Email..."
	if !sendEmail {
		verificationMessage = "Account created successfully, but verification email could not be sent."
	}

	if adminDowngraded {
		c.JSON(http.StatusCreated, gin.H{
			"message": fmt.Sprintf("%s Admin access revoked since Admin Password was incorrect.", verificationMessage),
			"user_id": user.UserID,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": verificationMessage,
		"user_id": user.UserID,
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

	if !user.IsVerified {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Please verify your email first",
			"user_id": user.UserID,
		})

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
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		MaxAge:   1800,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		Secure:   true,
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
		Name:     "access_token",
		Value:    "",
		Path:     "/",
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

func (uc *UserController) DeleteAccount(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req struct {
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	if strings.TrimSpace(req.Password) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password is required",
		})
		return
	}

	userId, err := utilities.GetUserIdfromContext(c)
	if err != nil {
		uc.dbLogger.Alerts("ACCOUNT", "Unable to fetch userId from context", gin.H{
			"endpoint": "/DeleteAccount",
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong!",
		})
		return
	}

	var user models.User

	usercollection := uc.db.Collection("users")

	err = usercollection.FindOne(ctx, bson.M{"user_id": userId}, options.FindOne().SetProjection(bson.M{"password": 1, "_id": 0})).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		uc.dbLogger.Alerts("ACCOUNT", "Unable to fetch user for account deletion", gin.H{
			"endpoint": "/DeleteAccount",
			"user_id":  userId,
			"error":    err.Error(),
		},
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong!",
		})
		return
	}

	err = utilities.ComparePassword(user.Password, req.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid credentials"})
		return
	}

	watchlistCollection := uc.db.Collection("watchlist")

	// Delete all watchlists owned by the user.
	_, err = watchlistCollection.DeleteMany(ctx, bson.M{"user_id": userId})
	if err != nil {
		uc.dbLogger.Alerts("ACCOUNT", "Unable to delete user watchlists", gin.H{
			"endpoint": "/DeleteAccount",
			"user_id":  userId,
			"error":    err.Error(),
		},
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to delete account data",
		})
		return
	}

	// Delete the user record last.
	result, err := usercollection.DeleteOne(ctx, bson.M{"user_id": userId})
	if err != nil {
		uc.dbLogger.Alerts("ACCOUNT", "Unable to delete user account",
			gin.H{"endpoint": "/DeleteAccount",
				"user_id": userId,
				"error":   err.Error(),
			},
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to delete account",
		})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User account not found",
		})
		return
	}

	// Remove refresh and access token cookie.
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	c.SetCookie("access_token", "", -1, "/", "", true, true)

	err = uc.otpService.DeleteVerificationOTP(ctx, userId)
	if err != nil {
		uc.dbLogger.Alerts("ERROR", "OTP Deletion failed", gin.H{
			"endpoint": "/DeleteAccount",
			"user_id":  userId,
			"error":    err.Error(),
		})
	}

	_ = uc.otpService.ClearAttempts(ctx, userId)
	if err != nil {
		uc.dbLogger.Alerts("ERROR", "failed in clearing attempts", gin.H{
			"endpoint": "/DeleteAccount",
			"user_id":  userId,
			"error":    err.Error(),
		})
	}

	uc.dbLogger.Log("ACCOUNT", "user deleted", gin.H{
			"endpoint": "/DeleteAccount",
			"user_id":  userId,
		},
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})

}

func (uc *UserController) VerifyEmail(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req models.VerifyEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	otp, err := uc.otpService.GetVerificationOTP(ctx, req.UserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "OTP expired or not found",
		})
		return
	}

	if otp != req.OTP {

		attempts, _ := uc.otpService.IncrementAttempts(ctx, req.UserID)

		if attempts >= 5 {

			uc.otpService.DeleteVerificationOTP(ctx, req.UserID)

			uc.otpService.ClearAttempts(ctx, req.UserID)

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Too many failed attempts. Request a new OTP.",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid OTP",
		})

		return
	}

	filter := bson.M{
		"user_id": req.UserID,
	}

	update := bson.M{
		"$set": bson.M{
			"is_verified": true,
			"updated_at":  time.Now(),
		},
	}

	usercollection := uc.db.Collection("users")

	_, err = usercollection.UpdateOne(ctx, filter, update)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to verify email",
		})

		return
	}

	uc.otpService.DeleteVerificationOTP(ctx, req.UserID)

	uc.otpService.ClearAttempts(ctx, req.UserID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
	})
}

func (uc *UserController) ResendVerification(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var req struct {
		UserId string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})

		return
	}

	cooldown, err := uc.otpService.IsResendCooldownActive(ctx, req.UserId)

	if err == nil && cooldown {

		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Please wait before requesting another OTP",
		})

		return
	}

	var user models.User

	usercollection := uc.db.Collection("users")

	err = usercollection.FindOne(ctx, bson.M{"user_id": req.UserId}, options.FindOne().SetProjection(bson.M{"_id": 0, "is_verified": 1, "email": 1})).Decode(&user)

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})

		return
	}

	if user.IsVerified {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already verified",
		})

		return
	}

	// Delete old OTPs
	err = uc.otpService.DeleteVerificationOTP(ctx, req.UserId)
	if err != nil {
		uc.dbLogger.Alerts("ERROR", "OTP Deletion failed", gin.H{
			"endpoint":   "/ResendVerification",
			"user_email": user.Email,
			"error":      err.Error(),
		})
	}

	sendEmail := true

	otp, err := uc.otpService.GenerateSecureOTP()
	if err != nil {
		uc.dbLogger.Alerts("ERROR", "OTP generation failed", gin.H{
			"endpoint":   "/ResendVerification",
			"user_email": user.Email,
			"error":      err.Error(),
		})
		sendEmail = false
	}

	if sendEmail {
		err = uc.otpService.SaveVerificationOTP(ctx, req.UserId, otp)

		if err != nil {
			uc.dbLogger.Alerts("ERROR", "Error in saving OTP", gin.H{
				"endpoint":   "/ResendVerification",
				"user_email": user.Email,
				"error":      err.Error(),
			})
			sendEmail = false
		}

		if sendEmail {
			err = services.SendVerificationOTP(user.Email, otp)
			if err != nil {
				uc.dbLogger.Alerts("ERROR", "Error sending email", gin.H{
					"endpoint":   "/ResendVerification",
					"user_email": user.Email,
					"error":      err.Error(),
				})
				sendEmail = false
			}
			uc.otpService.SetResendCooldown(ctx, req.UserId)
		}
	}

	if !sendEmail {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to send Verification code, try again.",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Verification code sent",
		})
	}
}

func (uc *UserController) ForgotPassword(c *gin.Context) {
	var req struct {
		EmailId string `json:"emailid"`
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	userCollection := uc.db.Collection("users")

	var user models.User

	err := userCollection.FindOne(ctx, bson.M{"email": req.EmailId}, options.FindOne().SetProjection(bson.M{
		"_id": 0, "user_id": 1,
	})).Decode(&user)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "If Account with provided email exists, reset link sent",
		})

		uc.dbLogger.Alerts("ERROR", "Error in fetching document count from DB", gin.H{
			"error":      err.Error(),
			"endpoint":   "/ForgotPassword",
			"user-email": req.EmailId,
		})

		return
	}

	resetPasswordToken := uuid.NewString()

	err = uc.otpService.SaveResetPasswordToken(ctx, user.UserID, resetPasswordToken)

	if err != nil {
		uc.dbLogger.Alerts("ERROR", "Error in storing reset password key in redis", gin.H{
			"error":      err.Error(),
			"endpoint":   "/ForgotPassword",
			"user-email": req.EmailId,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		return
	}

	// http://localhost:3000/reset-password?uniqueKey=abcdefg
	link := fmt.Sprintf("%s/reset-password?uniqueKey=%s", os.Getenv("FRONTEND_URL"), resetPasswordToken)

	err = services.SendPasswordResetOTP(req.EmailId, link)
	if err != nil {
		uc.dbLogger.Alerts("ERROR", "Error in sending password reset link", gin.H{
			"error":      err.Error(),
			"endpoint":   "/ForgotPassword",
			"user-email": req.EmailId,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to send link",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "If Account with provided email exists, reset link sent",
	})
}

func (uc *UserController) ResetPassword(c *gin.Context) {

	var req struct {
		Token           string `json:"token"`
		UpdatedPassword string `json:"password"`
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	userId, err := uc.otpService.GetResetPasswordToken(ctx, req.Token)
	if err != nil {
		uc.dbLogger.Alerts("ERROR", "Error in retrieving reset password token from redis", gin.H{
			"error":    err.Error(),
			"endpoint": "/ResetPassword",
			"token":    req.Token,
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or expired token",
		})
		return
	}

	updatedHashedPassword, err := utilities.HashPassword(req.UpdatedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
	}

	userCollection := uc.db.Collection("users")

	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
			"password":   updatedHashedPassword,
		},
	}

	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, update)

	if err != nil {
		uc.dbLogger.Alerts("ERROR", "Error in updating the password", gin.H{
			"error":    err.Error(),
			"endpoint": "/ResetPassword",
			"token":    req.Token,
		})

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	// delete token after use
	if err = uc.otpService.DeleteResetPasswordToken(ctx, req.Token); err != nil {
		uc.dbLogger.Alerts("ERROR", "Error in deleting reset password token from redis", gin.H{
			"error":    err.Error(),
			"endpoint": "/ResetPassword",
			"token":    req.Token,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password Reset successful"})

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

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    newToken,
		Path:     "/",
		MaxAge:   1800,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed"})
}

func (uc *UserController) GetProfile(c *gin.Context) {
	userId, err := utilities.GetUserIdfromContext(c)

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err != nil {
		uc.dbLogger.Alerts("ERROR", "Token update failed", gin.H{
			"endpoint": "/GetProfile",
			"userId":   userId,
			"error":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	usercollection := uc.db.Collection("users")

	var user models.User

	err = usercollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		Email:           user.Email,
		UserId:          user.UserID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Role:            user.Role,
		FavouriteGenres: user.FavouriteGenres,
	})

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
