package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yashbalyan08/rbac-vrv-go/controllers"
	"github.com/yashbalyan08/rbac-vrv-go/middleware"
)

func InitRoutes(router *gin.Engine) {
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	protected := router.Group("/admin")
	protected.Use(middleware.Authorize("Admin"))
	protected.GET("/dashboard", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to Admin Dashboard"})
	})
}
