package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	CORS() gin.HandlerFunc
}

type middleware struct {
	secretKey string
}

func NewMiddleware(secretKey string) Middleware {
	return &middleware{secretKey: secretKey}
}

// CORS ...
func (m *middleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
