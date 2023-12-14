// Package logmonitor provides logging utilities for a web application.
// It includes middleware that can be used with the Gin framework to log incoming
// HTTP requests and their response status, method, path, and the time taken to process.
//
// Copyright (c) 2023 H0llyW00dzZ
package logmonitor

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger returns a gin.HandlerFunc (middleware) that logs the details of HTTP requests.
// It records the status code, HTTP method, the path of the request, and the duration it took
// to process the request. This middleware is useful for monitoring and debugging the behavior
// of web services by providing insights into traffic patterns and potential bottlenecks.
//
// Example usage with the Gin framework:
//
//	r := gin.Default()
//	r.Use(logmonitor.RequestLogger())
//	// other routes and middleware
//	r.Run()
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer to track the duration of the request processing.
		start := time.Now()

		// Process the request by calling the next handler in the chain.
		c.Next()

		// Calculate the duration taken for the request to be processed.
		duration := time.Since(start)

		// Log the details of the request including the status code, method, path, and duration.
		fmt.Printf("Status: %d | Method: %s | Path: %s | Duration: %s\n",
			c.Writer.Status(),
			c.Request.Method,
			c.Request.URL.Path,
			duration,
		)
	}
}
