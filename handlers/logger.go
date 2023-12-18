package handlers

import (
	"github.com/H0llyW00dzZ/go-urlshortner/datastore"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
)

// Logger is a package-level variable to access the zap logger throughout the handlers package.
// It is intended to be used by other functions within the package for logging purposes.
var Logger *zap.Logger

// basePath is a package-level variable to store the base path for the handlers.
// It is set once during package initialization.
var basePath string

// internalSecretValue is a package-level variable that stores the secret value required by the InternalOnly middleware.
// It is set once during package initialization.
var internalSecretValue string

// SetLogger sets the logger instance for the package.
func SetLogger(logger *zap.Logger) {
	Logger = logger
}

// logAttemptToRetrieve logs an informational message indicating an attempt to retrieve the current URL by ID.
func logAttemptToRetrieve(id string) {
	logFields := createLogFields("retrieve", id) // Provide a default operation name
	Logger.Info(constant.AlertEmoji+"  "+constant.WarningEmoji+"  "+constant.InfoAttemptingToRetrieveTheCurrentURL, logFields...)
}

// logMismatchError logs an informational message indicating a mismatch error during URL update process.
func logMismatchError(id string) {
	logFields := createLogFields("mismatch_error", id) // Provide a default operation name
	Logger.Info(constant.GetBackEmoji+"  "+constant.UrlshortenerEmoji+"  "+constant.ErrorEmoji+"  "+constant.URLmismatchContextLog, logFields...)
}

// logAttemptToUpdate logs an informational message indicating an attempt to update a URL in the datastore.
func logAttemptToUpdate(id string) {
	logFields := createLogFields("update_attempt", id) // Provide a default operation name
	Logger.Info(constant.AlertEmoji+"  "+constant.WarningEmoji+"  "+datastore.InfoAttemptingToUpdateURLInDatastore, logFields...)
}

// logSuccessfulUpdate logs an informational message indicating a successful update of a URL in the datastore.
func logSuccessfulUpdate(id string) {
	logFields := createLogFields("successful_update", id) // Provide a default operation name
	Logger.Info(constant.UrlshortenerEmoji+"  "+constant.UpdateEmoji+"  "+constant.SuccessEmoji+"  "+datastore.InfoUpdateSuccessful, logFields...)
}

// LogInfo logs an informational message with given context fields.
func LogInfo(context string, fields ...zap.Field) {
	Logger.Info(context, fields...)
}

// LogError logs an error message with given context fields.
func LogError(context string, fields ...zap.Field) {
	Logger.Error(context, fields...)
}

// LogURLNotFound logs a "URL not found" error.
func LogURLNotFound(id string, err error) {
	logFields := createLogFieldsWithErr("getURL", id, err)
	LogInfo(constant.GetBackEmoji+"  "+constant.UrlshortenerEmoji+"  "+constant.URLnotfoundContextLog, logFields...)
}

// LogInternalError logs an internal server error.
func LogInternalError(context string, id string, err error) {
	logFields := createLogFieldsWithErr(context, id, err)
	LogError(constant.SosEmoji+"  "+constant.WarningEmoji+"  "+constant.FailedToGetURLContextLog, logFields...)
}

// LogURLRetrievalSuccess logs a successful URL retrieval.
func LogURLRetrievalSuccess(id string) {
	logFields := createLogFields("getURL", id) // Now correctly using two arguments
	Logger.Info(constant.UrlshortenerEmoji+"  "+constant.RedirectEmoji+"  "+constant.SuccessEmoji+"  "+constant.URLRetriveContextLog, logFields...)
}

// createLogFieldsWithErr is a helper to create log fields including an error.
func createLogFieldsWithErr(operation string, id string, err error) []zap.Field {
	return logmonitor.CreateLogFields(operation,
		logmonitor.WithComponent(constant.ComponentNoSQL),
		logmonitor.WithID(id),
		logmonitor.WithError(err),
	)
}

// createLogFields creates a slice of zap.Field with the operation and ID.
func createLogFields(operation, id string) []zap.Field {
	return []zap.Field{
		zap.String("operation", operation),
		zap.String("id", id),
	}
}
