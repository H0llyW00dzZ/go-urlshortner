package handlers

import (
	"github.com/H0llyW00dzZ/go-urlshortner/datastore"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"go.uber.org/zap"
)

// logAttemptToRetrieve logs an informational message indicating an attempt to retrieve the current URL by ID.
func logAttemptToRetrieve(id string) {
	logFields := createLogFields(id)
	logmonitor.Logger.Info(constant.AlertEmoji+"  "+constant.WarningEmoji+"  "+constant.InfoAttemptingToRetrieveTheCurrentURL, logFields...)
}

// logMismatchError logs an informational message indicating a mismatch error during URL update process.
func logMismatchError(id string) {
	logFields := createLogFields(id)
	logmonitor.Logger.Info(constant.GetBackEmoji+"  "+constant.UrlshortenerEmoji+"  "+constant.ErrorEmoji+"  "+constant.URLmismatchContextLog, logFields...)
}

// logAttemptToUpdate logs an informational message indicating an attempt to update a URL in the datastore.
func logAttemptToUpdate(id string) {
	logFields := createLogFields(id)
	logmonitor.Logger.Info(constant.AlertEmoji+"  "+constant.WarningEmoji+"  "+datastore.InfoAttemptingToUpdateURLInDatastore, logFields...)
}

// logSuccessfulUpdate logs an informational message indicating a successful update of a URL in the datastore.
func logSuccessfulUpdate(id string) {
	logFields := createLogFields(id)
	logmonitor.Logger.Info(constant.UrlshortenerEmoji+"  "+constant.UpdateEmoji+"  "+constant.SuccessEmoji+"  "+datastore.InfoUpdateSuccessful, logFields...)
}

// createLogFields generates a slice of zap.Field containing common log fields for the updateURL operation.
func createLogFields(id string) []zap.Field {
	return logmonitor.CreateLogFields("updateURL",
		logmonitor.WithComponent(constant.ComponentNoSQL),
		logmonitor.WithID(id),
	)
}
