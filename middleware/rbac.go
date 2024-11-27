package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yashbalyan08/rbac-vrv-go/utils"
)

func Authorize(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenStrings := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseJWT(tokenStrings)
		if err != nil || claims.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
