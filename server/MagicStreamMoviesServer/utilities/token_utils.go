package utilities

import (
	"errors"
	"os"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shivangsharma2925/MagicMovies/server/MagicStreamMoviesServer/models"
)

func GenerateAllTokens(email, firstname, lastname, role, userid string) (string, string, error) {
	claims := &models.SignedDetails{
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		Role:      role,
		UserId:    userid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicMovies",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		},
	}

	var secretKey string = os.Getenv("SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", "", err
	}

	refreshClaims := &models.SignedDetails{
		UserId: userid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicMovies",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	var refreshSecretKey string = os.Getenv("SECRET_REFRESH_KEY")

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshSignedToken, err := refreshToken.SignedString([]byte(refreshSecretKey))

	if err != nil {
		return "", "", err
	}

	return signedToken, refreshSignedToken, nil
}

func GetAccessToken(c *gin.Context) (string, error) {
	// authHeader := c.Request.Header.Get("Authorization")
	// if authHeader == "" {
	// 	return "", errors.New("Authorization header missing")
	// }

	// parts := strings.SplitN(authHeader, " ", 2)
	// if len(parts) != 2 {
	// 	return "", errors.New("invalid authorization header format")
	// }

	// // equalfold checks for Case-insensitive check, can be bearer, BEARER or Bearer
	// if !strings.EqualFold(parts[0], "Bearer") {
	// 	return "", errors.New("authorization type must be Bearer")
	// }

	// tokenString := strings.TrimSpace(parts[1])

	tokenString, err := c.Cookie("access_token")
	if err != nil {
		return "", err
	}

	if tokenString == "" {
		return "", errors.New("No token found")
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*models.SignedDetails, error) {
	claims := &models.SignedDetails{}

	var secretKey string = os.Getenv("SECRET_KEY")

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {

		//Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		logger.Warn("Expired token", "userId", claims.UserId)
		return nil, errors.New("Invalid or Expired token")
	}

	return claims, nil
}

func GetUserIdfromContext(c *gin.Context) (string, error) {
	userId, exists := c.Get("userId")
	if !exists {
		return "", errors.New("UserId does not exists in context")
	}

	id, ok := userId.(string)
	if !ok {
		return "", errors.New("Unable to fetch userId")
	}

	return id, nil
}

func GetRolefromContect(c *gin.Context) (string, error) {
	role, exists := c.Get("role")

	if !exists {
		return "", errors.New("Role not found in context")
	}

	roleString, ok := role.(string)

	if !ok {
		return "", errors.New("Unable to fetch role")
	}

	return roleString, nil
}

func ValidateRefreshToken(tokenString string) (*models.SignedRefreshDetails, error) {
	claims := &models.SignedRefreshDetails{}

	var secretKey string = os.Getenv("SECRET_REFRESH_KEY")

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		logger.Warn("Expired refresh token", "userId", claims.UserId)
		return nil, errors.New("refresh token has expired")
	}

	return claims, nil
}
