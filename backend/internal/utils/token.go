package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userId uint) (string, error) {
	// Generate a JWT token for the user
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	mySecret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(mySecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ExtractToken(authHeader string) (*jwt.Token, error) {
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}

func ExtractUserIdFromToken(token *jwt.Token) (uint, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["userId"] == nil {
		return 0, fmt.Errorf("invalid token claims")
	}
	userId, ok := claims["userId"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid userId type")
	}
	return uint(userId), nil
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the Authorization header is present
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header missing or invalid",
			})
			return
		}

		// Extract the token from the Authorization header
		token, err := ExtractToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		// Extract the user ID from the token claims
		userId, err := ExtractUserIdFromToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			return
		}

		// Store the user ID in the context for later use
		c.Set("userId", userId)
		c.Next()

	}
}
