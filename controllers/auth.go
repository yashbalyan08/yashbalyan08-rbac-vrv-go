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

	// Generate CSRF token (you can use a custom function here if needed)
	csrfToken := utils.GenerateCSRFToken()
	c.SetCookie(
		"auth_token",
		"Bearer "+token,
		3600*24,
		"",
		"",
		false,
		true,
	)
	c.SetCookie(
		"csrf_token",
		csrfToken,
		3600*24,
		"",
		"",
		false,
		true,
	)
	// Send the token and CSRF token in the response body for frontend usage
	c.JSON(http.StatusOK, gin.H{
		"message":    "Login successful",
		"jwt_token":  token,
		"csrf_token": csrfToken,
	})
}

// Logout handles user logout and token invalidation
func Logout(c *gin.Context) {
	// Invalidate the JWT token (if using cookies)
	c.SetCookie("auth_token", "", -1, "", "", false, true) // Set cookie expiration to -1 to clear it

	// Invalidate CSRF token
	c.SetCookie("csrf_token", "", -1, "", "", false, true) // Set cookie expiration to -1 to clear it

	// Optionally, remove the tokens from headers (if needed)
	// c.Header("Authorization", "")

	// Return a success message
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
