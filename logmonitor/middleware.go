package logmonitor

import (
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"github.com/gin-gonic/gin"
)

// Logger is a global variable to access the zap logger throughout the logmonitor package.
// It is initialized with a default configuration and can be replaced using SetLogger.
var Logger *zap.Logger

// SetLogger sets the logger instance for the package. This function allows for
// the replacement of the default logger with a customized one, enabling flexibility
// in logging configurations and output formats.
func SetLogger(logger *zap.Logger) {
	Logger = logger
}

// BadRequestError is a custom error type for bad requests. It includes a user-friendly
// message and the underlying error, providing context for logging and user feedback.
type BadRequestError struct {
	UserMessage string
	Err         error
}

// Error returns the message of the underlying error. This method allows BadRequestError
// to satisfy the error interface, making it compatible with Go's built-in error handling.
func (e *BadRequestError) Error() string {
	return e.Err.Error()
}

// NewBadRequestError creates a new instance of BadRequestError. This function is used
// to construct an error with a user-friendly message and an underlying error, which can
// be used to provide detailed error information while also giving a clear message to the end-user.
func NewBadRequestError(userMessage string, err error) *BadRequestError {
	return &BadRequestError{
		UserMessage: userMessage,
		Err:         err,
	}
}

// LogFieldOption defines a function signature for options that can be passed to createLogFields.
type LogFieldOption func() zap.Field

// CreateLogFields generates common log fields for use in various parts of the application.
func CreateLogFields(operation string, options ...LogFieldOption) []zap.Field {
	fields := []zap.Field{
		zap.String("operation", operation),
	}

	for _, opt := range options {
		fields = append(fields, opt())
	}

	return fields
}

// WithComponent returns a LogFieldOption that adds a 'component' field to the log.
func WithComponent(component string) LogFieldOption {
	return func() zap.Field {
		return zap.String("component", component)
	}
}

// WithID returns a LogFieldOption that adds an 'id' field to the log.
func WithID(id string) LogFieldOption {
	return func() zap.Field {
		return zap.String("id", id)
	}
}

// WithError returns a LogFieldOption that adds an 'error' field to the log.
func WithError(err error) LogFieldOption {
	return func() zap.Field {
		return zap.Error(err)
	}
}

// WithSignal returns a LogFieldOption that adds a 'signal' field to the log.
func WithSignal(signal os.Signal) LogFieldOption {
	return func() zap.Field {
		return zap.String(constant.ComponentMachineOperation, signal.String())
	}
}

// WithAnyZapField returns a LogFieldOption that allows direct use of zap.Field with CreateLogFields.
// This is useful for adding fields that are not covered by the other With* functions.
//
// Example:
//
//	logmonitor.CreateLogFields("operation hack the planet", WithAnyZapField(zap.Binary("H0llyW00dzZ", []byte("0x1337"))))
func WithAnyZapField(field zap.Field) LogFieldOption {
	return func() zap.Field {
		return field
	}
}

// RequestLogger returns a gin.HandlerFunc (middleware) that logs requests using zap.
// It captures key metrics for each HTTP or HTTPS request, including the status code, method,
// path, and processing duration, and outputs them in a structured format. This middleware
// enhances the observability of the application by providing detailed request logs, which
// are essential for monitoring and debugging.
//
// Upon receiving a request, it logs the following information:
//   - Machine Start Time (the local time when the request is received by the server)
//   - HTTP or HTTPS status code of the response
//   - HTTP or HTTPS method of the request
//   - Requested path
//   - Duration taken to process the request
//
// The logs are output in a structured format, making them easy to read and parse.
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer to track the duration of the request processing.
		start := time.Now()
		// Format the start time as a string in the desired format.
		startTimeFormatted := start.Format("2006/01/02 - 15:04:05")

		// Process the request by calling the next handler in the chain.
		c.Next()

		// Calculate the duration taken for the request to be processed.
		duration := time.Since(start)

		// Choose the emoji based on the HTTP status code.
		statusEmoji := constant.InfoEmoji
		if c.Writer.Status() >= 400 && c.Writer.Status() < 500 {
			statusEmoji = constant.WarningEmoji
		} else if c.Writer.Status() >= 500 {
			statusEmoji = constant.ErrorEmoji
		}

		// Log details of the request with zap, including the emoji.
		// Here we add the K8sEmoji to the log message.
		logger.Info(constant.K8sEmoji+"  "+statusEmoji+"  Request Details",
			zap.String("hostmachine_start_time", startTimeFormatted),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Duration("duration", duration),
		)
	}
}
