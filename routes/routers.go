package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yashbalyan08/rbac-vrv-go/controllers"
	"github.com/yashbalyan08/rbac-vrv-go/middleware"
)

func InitRoutes(router *gin.Engine) {
	// Route for login page
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	// Route for register page
	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})

	// Route for logged-in page
	router.GET("/logged-in", func(c *gin.Context) {
		c.HTML(http.StatusOK, "logged-in.html", nil)
	})

	// Route for registering a new user
	router.POST("/register", controllers.Register)

	// Route for user login
	router.POST("/login", controllers.Login)

	// Protected admin routes
	protected := router.Group("/admin")
	protected.Use(middleware.Authorize("Admin"))
	protected.GET("/dashboard", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to Admin Dashboard"})
	})
}
