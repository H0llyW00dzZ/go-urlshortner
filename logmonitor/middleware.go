// Package logmonitor provides logging utilities for a web application.
// It leverages the zap logging library to offer structured, leveled logging.
// The package is designed to integrate with the Gin web framework and includes
// middleware that logs incoming HTTP requests, including response status, method,
// path, and the time taken to process each request.
//
// The logger is initialized with a development-friendly configuration that outputs
// logs in a human-readable, color-coded format, suitable for development and debugging.
// The RequestLogger middleware can be easily added to a Gin router to enhance request
// logging with detailed information that can help in monitoring and troubleshooting.
//
// Example:
//
//	func main() {
//	    router := gin.Default()
//	    router.Use(logmonitor.RequestLogger())
//	    // ... other middlewares and routes ...
//	    router.Run(":8080")
//	}
//
// It is important to call the Logger.Sync() method when the application exits to
// flush any buffered log entries. This is typically done in the main function using
// a defer statement.
//
// Example:
//
//	func main() {
//	    defer func() {
//	        if err := logmonitor.Logger.Sync(); err != nil {
//	            fmt.Fprintf(os.Stderr, "Failed to flush log: %v\n", err)
//	        }
//	    }()
//	    // ... rest of the main function ...
//	}
//
// Copyright (c) 2023 H0llyW00dzZ
package logmonitor

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gin-gonic/gin"
)

// Logger is a global variable to access the zap logger throughout the logmonitor package.
var Logger *zap.Logger

func init() {
	// Initialize the zap logger with a development configuration.
	// This config is console-friendly and outputs logs in plaintext.
	config := zap.NewDevelopmentConfig()

	// Customize the level encoder to lowercase (info, warn, etc.)
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

	// Build the logger from the config and check for errors.
	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}
}

// RequestLogger returns a gin.HandlerFunc (middleware) that logs requests using zap.
// It is intended to be used as a middleware in a Gin router setup.
//
// Upon receiving a request, it logs the following information:
//   - HTTP status code of the response
//   - HTTP method of the request
//   - Requested path
//   - Duration taken to process the request
//
// The logs are output in a structured format, making them easy to read and parse.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer to track the duration of the request processing.
		start := time.Now()

		// Process the request by calling the next handler in the chain.
		c.Next()

		// Calculate the duration taken for the request to be processed.
		duration := time.Since(start)

		// Log details of the request with zap.
		Logger.Info("request details",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Duration("duration", duration),
		)
	}
}
