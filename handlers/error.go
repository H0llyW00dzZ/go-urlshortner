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
func handleDeletionError(c *gin.Context, err error) {
	id := c.Param(constant.HeaderID)
	// Use the centralized logging function from logmonitor package
	logFields := logmonitor.CreateLogFields("deleteURL",
		logmonitor.WithComponent(constant.ComponentGopher), // Use the constant for the component
		logmonitor.WithID(id),
		logmonitor.WithError(err),
	)
	// Fix internal issue now it's stable
	switch {
	case err == datastore.ErrNotFound:
		logmonitor.Logger.Info(constant.AlertEmoji+"  "+constant.WarningEmoji+"  "+constant.NoURLIDContextLog, logFields...)
		c.JSON(http.StatusNotFound, gin.H{
			constant.HeaderResponseError: constant.HeaderResponseIDandURLNotFound,
		})
	case strings.Contains(err.Error(), constant.URLmismatchContextLog):
		logmonitor.Logger.Info(constant.AlertEmoji+"  "+constant.WarningEmoji+"  "+constant.URLmismatchContextLog, logFields...)
		c.JSON(http.StatusBadRequest, gin.H{
			constant.HeaderResponseError: constant.URLmismatchContextLog,
		})
	case err.(*logmonitor.BadRequestError) != nil:
		badRequestErr := err.(*logmonitor.BadRequestError)
		logmonitor.Logger.Info(constant.AlertEmoji+"  "+constant.WarningEmoji+"  "+constant.HeaderResponseInvalidRequestJSONBinding, logFields...)
		c.JSON(http.StatusBadRequest, gin.H{
			constant.HeaderResponseError: badRequestErr.UserMessage,
		})
	default:
		logmonitor.Logger.Error(constant.SosEmoji+"  "+constant.WarningEmoji+"  "+constant.FailedToDeletedURLContextLog, logFields...)
		c.JSON(http.StatusInternalServerError, gin.H{
			constant.HeaderResponseError: constant.HeaderResponseInternalServerError,
		})
	}
}

// handleRetrievalError logs an error message for a failed retrieval attempt and returns a formatted error.
// If the error is a 'not found' error, it logs a specific message for that case.
func handleRetrievalError(err error, id string) error {
	logFields := createLogFields(id)
	if err == datastore.ErrNotFound {
		return datastore.ErrNotFound // Return the original error directly
	}
	logmonitor.Logger.Error(constant.FailedToRetriveURLContextLog+": "+err.Error(), logFields...)
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
