package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yashbalyan08/rbac-vrv-go/config"
	"github.com/yashbalyan08/rbac-vrv-go/models"
	"github.com/yashbalyan08/rbac-vrv-go/utils"
)

// Register handles the user registration process
func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	// Bind input JSON to struct and validate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := models.HashPassword(input.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password encryption failed"})
		return
	}

	// Create new user struct
	user := models.User{
		Username: input.Username,
		Password: hashedPassword,
		Role:     input.Role,
	}

	// Save the user in the database
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login handles user login and JWT generation
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Bind input JSON to struct and validate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	// Check if the user exists in the database
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify the password
	if !models.CheckPassword(user.Password, input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.Username, user.Role)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Set the JWT token as an HTTP-only cookie
	c.SetCookie(
		"auth_token", // Cookie name
		token,        // Token value
		3600,         // Max age in seconds (e.g., 3600 = 1 hour)
		"/",          // Path
		"",           // Domain (leave empty for default)
		false,        // Secure (set to true in production for HTTPS)
		true,         // HttpOnly (prevents access via JavaScript)
	)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func Logout(c *gin.Context) {
	// Set the auth_token cookie with a max age of -1 to delete it
	c.SetCookie(
		"auth_token", // Cookie name
		"",           // Empty value
		-1,           // Max age set to -1 to delete the cookie
		"/",          // Path
		"",           // Domain (leave empty for default)
		false,        // Secure (set to true in production for HTTPS)
		true,         // HttpOnly
	)

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
