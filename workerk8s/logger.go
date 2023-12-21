package workerk8s

import (
	"fmt"

	"go.uber.org/zap"
)

// Logger is a package-level variable to access the zap logger throughout the handlers package.
// It is intended to be used by other functions within the package for logging purposes.
var Logger *zap.Logger

// SetLogger sets the logger instance for the package.
func SetLogger(logger *zap.Logger) {
	Logger = logger
}

// logInfoWithEmoji logs an informational message with given emoji, context, and fields.
func logInfoWithEmoji(emoji string, context string, fields ...zap.Field) {
	Logger.Info(emoji+"  "+context, fields...)
}

// logErrorWithEmoji logs an error message with given emoji, context, and fields.
func logErrorWithEmoji(emoji string, context string, fields ...zap.Field) {
	Logger.Error(emoji+"  "+context, fields...)
}

// createLogFields creates a slice of zap.Field with the operation and additional info.
func createLogFields(operation string, namespace string, infos ...string) []zap.Field {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("namespace", namespace),
	}
	for i, info := range infos {
		fields = append(fields, zap.String(fmt.Sprintf("info%d", i+1), info))
	}
	return fields
}
