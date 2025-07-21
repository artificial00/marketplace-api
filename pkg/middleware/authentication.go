package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"marketplace-api/pkg/utils"
)

// AuthMiddleware проверяет JWT токен и добавляет пользователя в контекст
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.Unauthorized(c, "Bearer token required")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateToken(token, jwtSecret)
		if err != nil {
			utils.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_login", claims.Login)

		c.Next()
	}
}

// GetUserID извлекает ID пользователя из контекста
func GetUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(int)
	return id, ok
}
