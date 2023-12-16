package logmonitor

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gin-gonic/gin"
)

// Component constants for structured logging.
// This is used to identify the component that is logging the message.
const (
	ComponentNoSQL             = "datastore"
	ComponentCache             = "cache" // Currently unused.
	ComponentProjectIDENV      = "projectid"
	ComponentInternalSecretENV = "customsecretkey"
	ComponentMachineOperation  = "signal"
)

// Logger is a global variable to access the zap logger throughout the logmonitor package.
var Logger *zap.Logger

// SetLogger sets the logger instance for the package.
func SetLogger(logger *zap.Logger) {
	Logger = logger
}

// BadRequestError is a custom error type for bad requests.
type BadRequestError struct {
	UserMessage string
	Err         error
}

// LogFieldOption defines a function signature for options that can be passed to createLogFields.
type LogFieldOption func() zap.Field

func init() {
	// Initialize the zap logger with a development configuration.
	// This config is console-friendly and outputs logs in plaintext.
	config := zap.NewDevelopmentConfig()

	// Customize the logger configuration here if needed.
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Customize the level encoder to lowercase (info, warn, etc.)
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

	// Build the logger from the config and check for errors.
	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}
}

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

// WithComponent returns a LogFieldOption that adds an 'internal' field to the log.
func WithComponent(component string) LogFieldOption {
	return func() zap.Field {
		return zap.String("internal", component)
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
		return zap.String("notify", signal.String())
	}
}

// Error returns the message of the underlying error.
// This method allows BadRequestError to satisfy the error interface.
func (e *BadRequestError) Error() string {
	return e.Err.Error()
}

// NewBadRequestError creates a new instance of BadRequestError.
// This function is used to construct an error with a user-friendly message
// and an underlying error, which can be used to provide detailed error information
// while also giving a clear message to the end-user.
func NewBadRequestError(userMessage string, err error) *BadRequestError {
	return &BadRequestError{
		UserMessage: userMessage,
		Err:         err,
	}
}

// RequestLogger returns a gin.HandlerFunc (middleware) that logs requests using zap.
// It is intended to be used as a middleware in a Gin router setup.
//
// Upon receiving a request, it logs the following information:
//   - Machine Start Time (the local time when the request is received by the server)
//   - HTTP status code of the response
//   - HTTP method of the request
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

		// Log details of the request with zap.
		logger.Info("☸️  🗳️  Request Details",
			zap.String("machine_start_time", startTimeFormatted), // The local time when the request is received
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Duration("duration", duration),
		)
	}
}
