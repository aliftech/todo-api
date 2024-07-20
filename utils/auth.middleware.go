package utils

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aliftech/todo-api/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateAuth(c *gin.Context) {
	// Get the cookie from the request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized: No token provided.",
		})
		c.Abort()
		return
	}

	// Decode and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Retrieve the secret key from the environment variable
		secret := os.Getenv("jwt.secret")
		if secret == "" {
			return nil, fmt.Errorf("JWT secret key is not set in the environment variables")
		}

		return []byte(secret), nil
	})
	fmt.Print(err)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized: Invalid token.",
		})
		c.Abort()
		return
	}

	// Extract claims and validate them
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check token expiration
		if float64(time.Now().Unix()) > claims["expired"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Unauthorized: Token has expired.",
			})
			c.Abort()
			return
		}

		// Find the user with the token subject (sub)
		var user models.User
		if err := DB.First(&user, claims["sub"]).Error; err != nil || user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Unauthorized: User not found.",
			})
			c.Abort()
			return
		}

		// Attach the user ID to the request context
		c.Set("user", user.ID)

		// Continue to the next handler
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized: Invalid token claims.",
		})
		c.Abort()
		return
	}
}
