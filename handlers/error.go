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

// URLMismatchError represents an error for when the provided URL does not match
// the expected URL in the datastore. It embeds the error message to be returned.
type URLMismatchError struct {
	Message string
}

// Error returns the error message of a URLMismatchError.
// This method makes URLMismatchError satisfy the error interface.
func (e *URLMismatchError) Error() string {
	return e.Message
}

// isURLMismatchError checks if the provided error is of type URLMismatchError.
// It returns true if the error is a URLMismatchError, false otherwise.
// This is useful for type assertions where you need to identify if an error
// is specifically due to a URL mismatch.
func isURLMismatchError(err error) bool {
	_, ok := err.(*URLMismatchError) // Type assertion to check for URLMismatchError.
	return ok
}

// TODO: Refactor BadRequestError to be used consistently across the application.

// BadRequestError represents an error when the request made by the client
// contains bad syntax or cannot be fulfilled for some other reason.

type BadRequestError struct {
	Message string
}

// Error returns the error message of a BadRequestError.
// This method allows BadRequestError to satisfy the error interface,
// enabling it to be used like any other error.
// TODO:
// - Ensure BadRequestError is used in all handlers where bad requests need to be reported.
// - Update logging functions to handle BadRequestError specifically.
// - Review all BadRequestError occurrences to ensure the Message field is being used appropriately.

func (e *BadRequestError) Error() string {
	return e.Message
}

// handleUpdateError handles errors that occur during the URL update process.
func handleUpdateError(c *gin.Context, id string, err error) {
	logFields := logmonitor.CreateLogFields("editURL",
		logmonitor.WithComponent(constant.ComponentGopher),
		logmonitor.WithID(id),
		logmonitor.WithError(err),
	)

	switch {
	case err == datastore.ErrNotFound:
		logmonitor.Logger.Info(constant.GetBackEmoji+"  "+constant.UrlshortenerEmoji+"  "+constant.URLnotfoundContextLog, logFields...)
		c.JSON(http.StatusNotFound, gin.H{
			constant.HeaderResponseError: constant.URLnotfoundContextLog,
		})
	case strings.Contains(err.Error(), constant.URLmismatchContextLog):
		logmonitor.Logger.Info(constant.AlertEmoji+"  "+constant.WarningEmoji+"  "+constant.URLmismatchContextLog, logFields...)
		c.JSON(http.StatusBadRequest, gin.H{
			constant.HeaderResponseError: constant.URLmismatchContextLog,
		})
	default:
		// For other types of errors, respond with a 500 Internal Server Error.
		logmonitor.Logger.Error(constant.AlertEmoji+"  "+constant.WarningEmoji+"  "+constant.FailedToUpdateURLContextLog, logFields...)
		c.JSON(http.StatusInternalServerError, gin.H{
			constant.HeaderResponseError: constant.HeaderResponseInternalServerError,
		})
	}
}

// handleDeletionError handles errors that occur during the URL deletion process.
// Note this function `handleDeletionError` has maximum of 5 cyclomatic complexity so can't add another case here,
// because I don't have idea anymore for this function to reduce the cyclomatic complexity :v
func handleDeletionError(c *gin.Context, err error) {
	id := c.Param(constant.HeaderID)
	switch {
	case err == datastore.ErrNotFound:
		logNotFound(c, id) // Pass the context and id instead of err.Error()
		// This is a client error, so we should return a 400 status code.
	case isMismatchError(err):
		logMismatchError(c, id) // Pass the context and id instead of err.Error()
	case isBadRequestError(err):
		// This is also a client error bad request, so we should return a 400 status code.
		// Return a BadRequestError if JSON binding fails
		// Friendly error message for the user, and the original error for logging purposes
		logBadRequest(c, id) // Pass the context and id instead of err.Error()
	case isURLMismatchError(err): // This checks for the specific URL mismatch error
		logURLMismatchError(c, id, err)
	}
}

// SynclogError ensures that each error is logged only once.
func SynclogError(c *gin.Context, operation string, err error) {
	if err == nil || SyncerrorLogged(c) {
		return
	}

	// Pass the operation name to the logging function.
	SynclogSpecificError(c, operation, err)
	markErrorLogged(c)
}

// errorLogged checks if the error has already been logged.
func SyncerrorLogged(c *gin.Context) bool {
	_, logged := c.Get("errorLogged")
	return logged
}

// logSpecificError logs the error based on its type.
func SynclogSpecificError(c *gin.Context, operation string, err error) {
	switch err.(type) {
	case *BadRequestError:
		logBadRequest(c, operation) // Pass the operation name to the logging function.
	default:
		SynclogOtherError(c, operation, err) // Pass the operation name to the logging function.
	}
}

// logOtherError logs non-specific errors.
func SynclogOtherError(c *gin.Context, operation string, err error) {
	logFields := createLogFieldsWithErr(operation, "", err) // Use the operation name here.
	LogInfo(constant.ErrorEmoji+"  "+constant.UrlshortenerEmoji+"  "+constant.HeaderResponseInvalidRequest, logFields...)
}

// markErrorLogged marks the error as logged in the context.
func markErrorLogged(c *gin.Context) {
	c.Set("errorLogged", true)
}

// handleRetrievalError logs an error message for a failed retrieval attempt and returns a formatted error.
// If the error is a 'not found' error, it logs a specific message for that case.
func handleRetrievalError(err error, id string) error {
	// Add an operation name as the first argument to createLogFields.
	logFields := createLogFields("retrieveURL", id)
	if err == datastore.ErrNotFound {
		logmonitor.Logger.Info(constant.AlertEmoji+"  "+constant.WarningEmoji+"  "+constant.NoURLIDContextLog, logFields...)
		return datastore.ErrNotFound // Return the original error directly
	}
	logmonitor.Logger.Error(constant.SosEmoji+"  "+constant.WarningEmoji+"  "+constant.FailedToRetriveURLContextLog, logFields...)
	return err // Return the original error directly
}

// handleError logs the error and sends a JSON response with the error message and status code.
func handleError(c *gin.Context, message string, statusCode int, err error) {
	var emoji string

	// Check if the error is a BadRequestError and unwrap it if it is
	if badRequestErr, ok := err.(*logmonitor.BadRequestError); ok {
		message = badRequestErr.UserMessage
		statusCode = http.StatusBadRequest
		err = badRequestErr.Err
	}

	// Use different emojis based on the status code
	switch {
	case statusCode >= 500: // 5xx errors are still logged as errors
		emoji = constant.ErrorEmoji
		Logger.Error(emoji+"  "+message, zap.Error(err))
	}

	c.AbortWithStatusJSON(statusCode, gin.H{
		constant.HeaderResponseError: message,
	})
}
