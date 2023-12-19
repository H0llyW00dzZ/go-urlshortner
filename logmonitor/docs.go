// Package logmonitor provides structured, leveled logging utilities designed for web applications.
// It utilizes the zap logging library to facilitate structured logging and integrates seamlessly
// with the Gin web framework. The package includes middleware for logging HTTP or HTTPS requests, capturing
// key information such as response status, method, path, and processing time for each request.
//
// # Available Functions, Variables, and Types:
//
// # Variables:
//   - Logger: A global *zap.Logger variable for logging throughout the application.
//
// # Functions:
//   - SetLogger(logger *zap.Logger): Sets the global Logger variable to a specified zap.Logger.
//   - NewBadRequestError(userMessage string, err error): Creates a new BadRequestError instance.
//   - CreateLogFields(operation string, options ...LogFieldOption): Generates common log fields.
//   - WithComponent(component string): Returns a LogFieldOption that adds a 'component' field.
//   - WithID(id string): Returns a LogFieldOption that adds an 'id' field.
//   - WithError(err error): Returns a LogFieldOption that adds an 'error' field.
//   - WithSignal(signal os.Signal): Returns a LogFieldOption that adds a 'signal' field.
//   - RequestLogger(logger *zap.Logger): Gin middleware that logs HTTP or HTTPS requests.
//
// # Types:
//   - BadRequestError: Custom error type with a user-friendly message and an underlying error.
//   - LogFieldOption: Function signature for options to create log fields.
//
// # Constants:
//   - The package may also define constants for log categorization (not shown in the provided code).
//
// # Usage example:
//
//	func main() {
//	    router := gin.Default()
//	    router.Use(logmonitor.RequestLogger(logmonitor.Logger))
//	    // ... additional middleware and route setup ...
//	    router.Run(":8080")
//	}
//
// Flushing Logs:
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
// The RequestLogger middleware logs vital request details and should be included as part of the
// Gin router setup to capture request metrics in a structured log format.
//
// Copyright (c) 2023 by H0llyW00dzZ
package logmonitor
