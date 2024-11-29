package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yashbalyan08/rbac-vrv-go/utils"
)

// AuthMiddleware ensures the user is authenticated using JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the Authorization header
		authHeader, _ := c.Cookie("auth_token")
		log.Println("Header: ", authHeader)
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		// Extract the token and validate it
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		log.Println(tokenString)
		_, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Proceed to the next handler
		c.Next()
	}
}

// AuthorizeRoles ensures the user has one of the allowed roles using JWT
func AuthorizeRoles(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the Authorization header
		t, err := c.Cookie("auth_token")
		log.Println(t)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Extract the token and validate it
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		// Check if the user's role is in the list of allowed roles
		for _, role := range allowedRoles {
			if claims.Role == role {
				// Proceed to the next handler if the role is authorized
				c.Next()
				return
			}
		}

		// Deny access if no roles match
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		c.Abort()
	}
}
