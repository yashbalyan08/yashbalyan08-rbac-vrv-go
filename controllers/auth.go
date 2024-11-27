package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yashbalyan08/rbac-vrv-go/config"
	"github.com/yashbalyan08/rbac-vrv-go/models"
	"github.com/yashbalyan08/rbac-vrv-go/utils"
)

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedpassword, err := models.HashPassword(input.Password)
	if err != nil {
		log.Fatalf("Error hashing: %v", err)
	}
	user := models.User{
		Username: input.Username,
		Password: hashedpassword,
		Role:     input.Role,
	}
	if err = config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Registered Successfully"})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
		return
	}

	if !models.CheckPassword(user.Password, input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
