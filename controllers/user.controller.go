package controllers

import (
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
		Email    string
		Fullname string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Failed to read request body",
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
			"message": "Failed to read request body",
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

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     user.ID,
		"expired": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("jwt.secret")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create token",
			"data":    nil,
		})
		return
	}

	// Set the cookie with the JWT token
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Login Success.",
		"data":    nil,
	})
}
