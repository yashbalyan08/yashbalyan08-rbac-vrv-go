package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yashbalyan08/rbac-vrv-go/controllers"
	"github.com/yashbalyan08/rbac-vrv-go/middleware"
)

func InitRoutes(router *gin.Engine) {
	// Login and Register routes
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	router.GET("/register", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	// Profile route (requires authentication)
	router.GET("/profile", middleware.AuthMiddleware(), middleware.AuthorizeRoles([]string{"User"}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "User profile data"})
	})

	// Dashboard route (requires authentication)
	router.GET("/dashboard", middleware.AuthMiddleware(), middleware.AuthorizeRoles([]string{"Admin"}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Dashboard for authenticated users"})
	})

	// Logged-in route
	router.GET("/logged-in", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "logged-in.html", nil)
	})

	router.POST("/logout", controllers.Logout)

}
