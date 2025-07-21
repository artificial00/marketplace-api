package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware настраивает CORS заголовки
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, Authorization, X-Requested-With")

		c.Header("Access-Control-Allow-Credentials", "true")

		c.Header("Access-Control-Max-Age", "86400") // 24 часа

		if c.Request.Method == "OPTIONS" {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
			c.Abort()
			return
		}

		c.Next()
	}
}
