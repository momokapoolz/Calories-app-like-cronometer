package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs HTTP requests with time taken to process
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Process request
		c.Next()
		
		// Log the request details
		log.Printf(
			"%s %s %s %s",
			c.Request.Method,
			c.Request.RequestURI,
			c.Request.RemoteAddr,
			time.Since(start),
		)
	}
} 