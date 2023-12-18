package handlers

import (
	"net/http"
	"strings"

	"github.com/H0llyW00dzZ/go-urlshortner/datastore"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"github.com/gin-gonic/gin"
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

// logInfoWithEmoji logs an informational message with given emoji, context, and fields.
func logInfoWithEmoji(emoji string, context string, fields ...zap.Field) {
	Logger.Info(emoji+"  "+context, fields...)
}

// logErrorWithEmoji logs an error message with given emoji, context, and fields.
func logErrorWithEmoji(emoji string, context string, fields ...zap.Field) {
	Logger.Error(emoji+"  "+context, fields...)
}

// logAttemptToRetrieve logs an informational message indicating an attempt to retrieve the current URL by ID.
func logAttemptToRetrieve(id string) {
	logFields := createLogFields("retrieve", id) // Provide a default operation name
	Logger.Info(constant.AlertEmoji+"  "+constant.WarningEmoji+"  "+constant.InfoAttemptingToRetrieveTheCurrentURL, logFields...)
}

// LogMismatchError logs a message indicating that there is a mismatch error.
func LogMismatchError(id string) {
	logFields := createLogFields("mismatch_error", id)
	Logger.Info(constant.ErrorEmoji+" "+constant.UrlshortenerEmoji+" "+constant.HeaderResponseInvalidRequestPayload, logFields...)
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
	fields := createLogFieldsWithErr("getURL", id, err)
	logInfoWithEmoji(constant.GetBackEmoji+" "+constant.UrlshortenerEmoji, constant.URLnotfoundContextLog, fields...)
}

// LogInternalError logs an internal server error.
func LogInternalError(context string, id string, err error) {
	fields := createLogFieldsWithErr(context, id, err)
	logErrorWithEmoji(constant.SosEmoji+" "+constant.WarningEmoji, constant.FailedToGetURLContextLog, fields...)
}

// LogURLRetrievalSuccess logs a successful URL retrieval.
func LogURLRetrievalSuccess(id string) {
	logFields := createLogFields("getURL", id) // Now correctly using two arguments
	Logger.Info(constant.UrlshortenerEmoji+"  "+constant.RedirectEmoji+"  "+constant.SuccessEmoji+"  "+constant.URLRetriveContextLog, logFields...)
}

// LogInvalidURLFormat logs a message indicating that the URL format is invalid.
func LogInvalidURLFormat(url string) {
	Logger.Info(constant.ErrorEmoji+"  "+constant.UrlshortenerEmoji+"  "+constant.HeaderResponseInvalidURLFormat,
		zap.String("url", url),
	)
}

// LogBadRequestError logs a message indicating a bad request error.
func LogBadRequestError(context string, err error) {
	fields := createLogFieldsWithErr(context, "", err) // Assuming an empty ID for general bad requests
	logInfoWithEmoji(constant.ErrorEmoji+" "+constant.UrlshortenerEmoji, constant.HeaderResponseInvalidRequestPayload, fields...)
}

// LogURLShortened logs a message indicating that a URL has been successfully shortened.
func LogURLShortened(id string) {
	logFields := createLogFields("shorten_url", id)
	Logger.Info(constant.UrlshortenerEmoji+"  "+constant.SuccessEmoji+"  "+constant.URLShorteneredContextLog, logFields...)
}

// LogDeletionError logs a message indicating that there was an error during deletion.
func LogDeletionError(id string, err error) {
	logFields := createLogFieldsWithErr("delete_url", id, err)
	LogInfo(constant.ErrorEmoji+"  "+constant.UrlshortenerEmoji+"  "+constant.ErrorDuringDeletionContextLog, logFields...)
}

// LogURLDeletionSuccess logs a message indicating that a URL has been successfully deleted.
func LogURLDeletionSuccess(id string) {
	logFields := createLogFields("delete_url", id)
	Logger.Info(constant.UrlshortenerEmoji+"  "+constant.SuccessEmoji+"  "+constant.URLDeletedSuccessfullyContextLog, logFields...)
}

// Use the centralized logging function from logmonitor package
func createDeletionLogFields(id string, err error) []zap.Field {
	return logmonitor.CreateLogFields("deleteURL",
		logmonitor.WithComponent(constant.ComponentGopher),
		logmonitor.WithID(id),
		logmonitor.WithError(err),
	)
}

// logNotFound handles logging and response for a "not found" situation.
func logNotFound(c *gin.Context, id string) {
	fields := createLogFields("deleteURL", id)
	logInfoWithEmoji(constant.AlertEmoji+"  "+constant.WarningEmoji, constant.NoURLIDContextLog, fields...)
	c.JSON(http.StatusNotFound, gin.H{
		constant.HeaderResponseError: constant.HeaderResponseIDandURLNotFound,
	})
}

// isMismatchError checks if the error is a "mismatch error" situation.
func isMismatchError(err error) bool {
	return strings.Contains(err.Error(), constant.PathIDandPayloadIDDoesnotMatchContextLog)
}

// logMismatchError handles logging and response for a "mismatch error" situation.
func logMismatchError(c *gin.Context, id string) {
	fields := createLogFields("deleteURL", id)
	logInfoWithEmoji(constant.ErrorEmoji+"  "+constant.WarningEmoji, constant.URLmismatchContextLog, fields...)
	c.JSON(http.StatusBadRequest, gin.H{
		constant.HeaderResponseError: constant.PathIDandPayloadIDDoesnotMatchContextLog,
	})
}

// logURLMismatchError logs the URL mismatch error with appropriate emojis and sends a
// 400 Bad Request response with the URL mismatch error message.
func logURLMismatchError(c *gin.Context, id string, err error) {
	// Create log fields to include additional metadata in the log entry.
	fields := createLogFields("url_mismatch_error", id)

	// Log the error with an information level log entry, including emojis for visibility.
	logInfoWithEmoji(constant.ErrorEmoji+"  "+constant.WarningEmoji, constant.URLmismatchContextLog, fields...)

	// Respond to the client with a 400 Bad Request status code and include the error message.
	// This indicates that the server cannot process the request due to a client error (mismatched URL).
	c.JSON(http.StatusBadRequest, gin.H{
		constant.HeaderResponseError: constant.URLmismatchContextLog,
	})
}

// isBadRequestError checks if the error is a "bad request" situation.
func isBadRequestError(err error) bool {
	return strings.Contains(err.Error(), constant.HeaderResponseInvalidRequestPayload)
}

// logBadRequest handles logging and response for a "bad request" situation.
func logBadRequest(c *gin.Context, id string) {
	fields := createLogFields("deleteURL", id)
	logInfoWithEmoji(constant.ErrorEmoji+"  "+constant.WarningEmoji, constant.HeaderResponseInvalidRequestJSONBinding, fields...)
	c.JSON(http.StatusBadRequest, gin.H{
		constant.HeaderResponseError: constant.HeaderResponseInvalidRequestPayload,
	})
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
