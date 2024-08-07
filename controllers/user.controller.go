package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aliftech/todo-api/models"
	"github.com/aliftech/todo-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(c *gin.Context) {
	// Get email and password of req body
	var body struct {
		Email    string `form:"email" json:"email" binding:"required,email"`
		Fullname string `form:"fullname" json:"fullname" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request value.",
			"data":    nil,
		})

		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Failed to hash password",
			"data":    nil,
		})

		return
	}

	// Create the user
	user := models.User{Email: body.Email, Password: string(hash)}

	result := utils.DB.Create(&user) // pass pointer of data to Create

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Failed to create new user",
			"data":    nil,
		})

		return
	}

	// Show response
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "New user have been created",
		"data":    nil,
	})
}

func Login(c *gin.Context) {
	// Get email and password from request body
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request value.",
			"data":    nil,
		})
		return
	}

	// Lookup the requested user
	var user models.User
	if result := utils.DB.Where("email = ?", body.Email).First(&user); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid email or password",
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Database error",
				"data":    nil,
			})
		}
		return
	}

	// Compare the provided password with the saved hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid email or password",
			"data":    nil,
		})
		return
	}

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     user.ID,
		"expired": time.Now().Add(time.Hour * 1).Unix(), // 1 hour
	})

	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("jwt.secret")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create access token",
			"data":    nil,
		})
		return
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     user.ID,
		"expired": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("jwt.secret")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create refresh token",
			"data":    nil,
		})
		return
	}

	// Return the tokens in the response body
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Login Success.",
		"data": gin.H{
			"access_token":  accessTokenString,
			"refresh_token": refreshTokenString,
		},
	})
}

func Refresh(c *gin.Context) {
	// Get the refresh token from the request body
	var body struct {
		Token string `binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request value.",
			"data":    nil,
		})
		return
	}

	// Parse and validate the refresh token
	refreshToken, err := jwt.Parse(body.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		secret := os.Getenv("jwt.secret")
		if secret == "" {
			return nil, fmt.Errorf("JWT secret key is not set in the environment variables")
		}

		return []byte(secret), nil
	})

	if err != nil || !refreshToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized: Invalid refresh token.",
			"data":    nil,
		})
		return
	}

	// Extract claims and validate them
	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		// Check token expiration
		exp, ok := claims["expired"].(float64)
		if !ok || float64(time.Now().Unix()) > exp {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Unauthorized: Refresh token has expired.",
				"data":    nil,
			})
			return
		}

		// Generate new access token
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":     claims["sub"],
			"expired": time.Now().Add(time.Hour * 1).Unix(), // 1 hour
		})

		accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("jwt.secret")))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to create access token",
				"data":    nil,
			})
			return
		}

		// Generate new refresh token
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":     claims["sub"],
			"expired": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
		})

		refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("jwt.secret")))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to create refresh token",
				"data":    nil,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Token refreshed successfully.",
			"data": gin.H{
				"access_token":  accessTokenString,
				"refresh_token": refreshTokenString,
			},
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized: Invalid token claims.",
			"data":    nil,
		})
		return
	}
}
