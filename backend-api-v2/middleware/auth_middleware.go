package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const MockToken = "mock-jwt-token"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "4001",
				"message": "Unauthorized: Missing token",
				"data":    nil,
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "4001",
				"message": "Unauthorized: Invalid token format",
				"data":    nil,
			})
			c.Abort()
			return
		}

		token := parts[1]
		if token != MockToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "4001",
				"message": "Unauthorized: Invalid token",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Set("username", "intern_user")
		c.Next()
	}
}
