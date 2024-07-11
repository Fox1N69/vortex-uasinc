package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

type Middleware interface {
	CORS() gin.HandlerFunc
	RPSLimit(rps int) gin.HandlerFunc
}

type middleware struct {
	secretKey string
}

func NewMiddleware(secretKey string) Middleware {
	return &middleware{secretKey: secretKey}
}

// CORS returns a middleware handler that adds CORS headers to the response.
//
// It sets the following headers:
//   - Access-Control-Allow-Origin: *
//   - Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
//   - Access-Control-Allow-Headers: Origin, Content-Type, Authorization
//   - Access-Control-Expose-Headers: Content-Length
//   - Access-Control-Allow-Credentials: true
//
// If the incoming request method is OPTIONS, it responds with HTTP status
// 204 (No Content) and aborts further processing.
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

// RPSLimit returns a middleware handler that limits the requests per second (RPS).
//
// It uses a rate limiter initialized with the provided RPS value. For each request,
// it logs the time elapsed since the previous request using the rate limiter.
func (m *middleware) RPSLimit(rps int) gin.HandlerFunc {
	limit := ratelimit.New(rps)
	prev := time.Now()

	return func(c *gin.Context) {
		now := limit.Take()
		log.Printf("Time since last request: %v", now.Sub(prev))
		prev = now
	}
}
