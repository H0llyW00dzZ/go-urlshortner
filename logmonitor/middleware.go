package logmonitor

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger returns a gin.HandlerFunc (middleware) that logs requests.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate resolution time
		duration := time.Since(start)

		// Log details of the request
		fmt.Printf("Status: %d | Method: %s | Path: %s | Duration: %s\n",
			c.Writer.Status(),
			c.Request.Method,
			c.Request.URL.Path,
			duration,
		)
	}
}
