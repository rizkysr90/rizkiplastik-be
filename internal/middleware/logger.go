package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger middleware for request logging
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate request processing time
		latency := time.Since(startTime)

		// Get request details
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// Log request details
		log.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %s",
			startTime.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
