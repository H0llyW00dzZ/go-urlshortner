// Package logmonitor provides structured, leveled logging utilities designed for web applications.
// It utilizes the zap logging library to facilitate structured logging and integrates seamlessly
// with the Gin web framework. The package includes middleware for logging HTTP requests, capturing
// key information such as response status, method, path, and processing time for each request.
//
// The logger is configured with a development-friendly setup that outputs logs in a color-coded,
// human-readable format, which is ideal for development and debugging purposes. The RequestLogger
// middleware can be readily applied to a Gin engine to augment request logging with granular details,
// aiding in application monitoring and troubleshooting.
//
// Usage example:
//
//	func main() {
//	    router := gin.Default()
//	    router.Use(logmonitor.RequestLogger(logmonitor.Logger))
//	    // ... additional middleware and route setup ...
//	    router.Run(":8080")
//	}
//
// It is crucial to flush any buffered log entries upon application termination to ensure all logs
// are committed to their intended destination. This is accomplished by invoking the Logger.Sync()
// method, typically in the main function using defer to guarantee execution even during an
// unexpected exit.
//
// Flush logs on application exit example:
//
//	func main() {
//	    defer func() {
//	        if err := logmonitor.Logger.Sync(); err != nil {
//	            fmt.Fprintf(os.Stderr, "Failed to flush logs: %v\n", err)
//	        }
//	    }()
//	    // ... remainder of the main function ...
//	}
//
// The package also defines constants for various components to categorize logs, allowing for
// filtering and analysis based on specific parts of the application. These constants are utilized
// when creating log fields, ensuring consistent identification across logs.
//
// This package's global Logger variable provides access to the configured zap logger for use
// throughout the application. SetLogger allows for the replacement of the default logger with
// a customized one if necessary.
//
// Additional types and functions, such as BadRequestError and associated constructors, offer
// convenience in generating structured logs with common fields and handling specific error scenarios.
//
// The RequestLogger middleware logs vital request details and should be included as part of the
// Gin router setup to capture request metrics in a structured log format.
//
// Copyright (c) 2023 H0llyW00dzZ
package logmonitor
