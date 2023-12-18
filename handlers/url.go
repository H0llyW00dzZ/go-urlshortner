package handlers

import (
	"fmt"
	"net/http"

	"github.com/H0llyW00dzZ/go-urlshortner/datastore"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// getURLHandlerGin returns a Gin handler function that retrieves and redirects to the original
// URL based on a short identifier provided in the request path. If the identifier is not found
// or an error occurs, the handler responds with the appropriate HTTP status code and error message.
func getURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param(constant.HeaderID)

		// Assuming datastore.GetURL is a function that correctly handles datastore operations.
		url, err := datastore.GetURL(c, dsClient, id)
		// Declare logFields here so it's accessible throughout the function scope
		logFields := logmonitor.CreateLogFields("getURL",
			logmonitor.WithComponent(constant.ComponentNoSQL), // Use the constant for the component
			logmonitor.WithID(id),
			logmonitor.WithError(err), // Include the error here, but it will be nil if there's no error
		)

		if err != nil {
			if err == datastore.ErrNotFound {
				logmonitor.Logger.Info(constant.GetBackEmoji+"  "+constant.UrlshortenerEmoji+"  "+constant.URLnotfoundContextLog, logFields...)
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					constant.HeaderResponseError: constant.URLnotfoundContextLog,
				})
			} else {
				logmonitor.Logger.Error(constant.SosEmoji+"  "+constant.WarningEmoji+"  "+constant.FailedToGetURLContextLog, logFields...)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					constant.HeaderResponseError: constant.HeaderResponseInternalServerError,
				})
			}
			return
		}

		// Check if URL is nil after the GetURL call
		if url == nil {
			// Use the logmonitor's logging function for consistency
			logmonitor.Logger.Error(constant.SosEmoji+"  "+constant.WarningEmoji+"  "+constant.URLisNilContextLog, logFields...)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				constant.HeaderResponseError: constant.HeaderResponseInternalServerError,
			})
			return
		}

		// If there's no error and you're logging a successful retrieval, use the same logFields
		logmonitor.Logger.Info(constant.UrlshortenerEmoji+"  "+constant.RedirectEmoji+"  "+constant.SuccessEmoji+"  "+constant.URLRetriveContextLog, logFields...)
		c.Redirect(http.StatusFound, url.Original)
	}
}

// postURLHandlerGin returns a Gin handler function that handles the creation of a new shortened
// URL. It expects a JSON payload with the original URL, generates a short identifier, and stores
// the mapping in Google Cloud Datastore. If successful, it returns the generated identifier and
// the shortened URL; otherwise, it responds with an error.
func postURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract and validate the original URL from the request body.
		url, err := extractURL(c)
		if err != nil {
			handleError(c, constant.HeaderResponseInvalidRequestPayload, http.StatusBadRequest, err)
			return
		}

		// Generate a short identifier for the URL.
		id, err := generateShortID(c.Request.Context(), dsClient) // Pass the request context and datastore client
		if err != nil {
			handleError(c, constant.HeaderResponseFailedtoGenerateID, http.StatusInternalServerError, err)
			return
		}

		// Save the URL with the generated identifier into the datastore.
		if err := saveURL(c, dsClient, id, url); err != nil {
			handleError(c, constant.HeaderResponseFailedtoSaveURL, http.StatusInternalServerError, err)
			return
		}

		logFields := logmonitor.CreateLogFields("postURL",
			logmonitor.WithComponent(constant.ComponentNoSQL), // Use the constant for the component
			logmonitor.WithID(id),
		)

		logmonitor.Logger.Info(constant.UrlshortenerEmoji+"  "+constant.SuccessEmoji+"  "+constant.URLShorteneredContextLog, logFields...)

		// Construct the full shortened URL and return it in the response.
		fullShortenedURL := constructFullShortenedURL(c, id)
		c.JSON(http.StatusOK, gin.H{
			constant.HeaderID: id, constant.HeaderResponseshortened_url: fullShortenedURL,
		})
	}
}

// editURLHandlerGin returns a Gin handler function that handles the updating of an existing shortened URL.
func editURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		pathID, req, err := validateUpdateRequest(c)
		if err != nil {
			// Handle the error, including the case where the JSON fields don't match
			handleError(c, constant.HeaderResponseInvalidRequestPayload, http.StatusBadRequest, err)
			return
		}

		err = updateURL(c, dsClient, pathID, req)
		if err != nil {
			handleUpdateError(c, pathID, err)
			return
		}

		respondWithUpdatedURL(c, pathID)
	}
}

// validateUpdateRequest validates the update request and extracts the path ID and request payload.
func validateUpdateRequest(c *gin.Context) (pathID string, req UpdateURLPayload, err error) {
	pathID = c.Param(constant.HeaderID)
	if err := c.ShouldBindJSON(&req); err != nil {
		logFields := logmonitor.CreateLogFields("validateUpdateRequest",
			logmonitor.WithComponent(constant.ComponentGopher),
			logmonitor.WithID(pathID),
			logmonitor.WithError(err),
		)
		// Return a BadRequestError if JSON binding fails
		// Friendly error message for the user, and the original error for logging purposes
		// Log the error with structured logging
		logmonitor.Logger.Info(constant.ErrorEmoji+"  "+constant.UrlshortenerEmoji+"  "+constant.HeaderResponseInvalidRequestJSONBinding, logFields...)

		// Return a BadRequestError with the actual error
		return "", req, err
	}

	// Additional validation for the ID in the URL and the ID in the payload
	if pathID != req.ID {
		err := fmt.Errorf(constant.PathIDandPayloadIDDoesnotMatchContextLog)
		logFields := logmonitor.CreateLogFields("validateUpdateRequest",
			logmonitor.WithComponent(constant.ComponentGopher),
			logmonitor.WithID(pathID),
			logmonitor.WithError(err),
		)
		logmonitor.Logger.Info(constant.ErrorEmoji+"  "+constant.UrlshortenerEmoji+"  "+constant.HeaderResponseInvalidRequestPayload, logFields...)
		return "", req, err
	}

	return pathID, req, nil
}

// updateURL retrieves the current URL, verifies it against the provided old URL, and updates it with the new URL.
// It returns an error with a message suitable for HTTP response if any step fails.
func updateURL(c *gin.Context, dsClient *datastore.Client, id string, req UpdateURLPayload) error {
	logAttemptToRetrieve(id)

	currentURL, err := datastore.GetURL(c, dsClient, id)
	if err != nil {
		return handleRetrievalError(err, id)
	}

	if currentURL.Original != req.OldURL {
		logMismatchError(id)
		return fmt.Errorf(constant.URLmismatchContextLog)
	}

	logAttemptToUpdate(id)

	// Update the URL in the datastore with the new URL.
	if err := datastore.UpdateURL(c, dsClient, id, req.NewURL); err != nil {
		return err // Simply return the error
	}

	logSuccessfulUpdate(id)

	return nil
}

// respondWithUpdatedURL constructs and sends a JSON response with the updated URL information.
func respondWithUpdatedURL(c *gin.Context, id string) {
	fullShortenedURL := constructFullShortenedURL(c, id)
	c.JSON(http.StatusOK, gin.H{
		constant.HeaderID:                    id,
		constant.HeaderResponseshortened_url: fullShortenedURL,
		constant.HeaderResponseStatus:        constant.HeaderResponseURlUpdated,
	})
}

// extractURL extracts the original URL from the JSON payload in the request.
func extractURL(c *gin.Context) (string, error) {
	var req CreateURLPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		Logger.Info(constant.ErrorEmoji+"  "+constant.UrlshortenerEmoji+"  "+constant.HeaderResponseInvalidRequestJSONBinding, zap.Error(err))
		return "", err
	}

	// Check if the URL is in a valid format.
	if req.URL == "" || !isValidURL(req.URL) {
		Logger.Info(constant.ErrorEmoji+"  "+constant.UrlshortenerEmoji+"  "+constant.HeaderResponseInvalidURLFormat, zap.String("url", req.URL))
		return "", fmt.Errorf(constant.HeaderResponseInvalidURLFormat)
	}

	return req.URL, nil
}

// deleteURLHandlerGin returns a Gin handler function that handles the deletion of an existing shortened URL.
func deleteURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use the centralized logging function from logmonitor package
		logFields := logmonitor.CreateLogFields("deleteURL",
			logmonitor.WithComponent(constant.ComponentNoSQL), // Use the constant for the component
			logmonitor.WithID(c.Param(constant.HeaderID)),
		)
		if err := validateAndDeleteURL(c, dsClient); err != nil {
			handleDeletionError(c, err)
		} else {
			logmonitor.Logger.Info(constant.DeleteEmoji+"  "+constant.UrlshortenerEmoji+"  "+constant.SuccessEmoji+"  "+constant.HeaderResponseURLDeleted, logFields...)
			c.JSON(http.StatusOK, gin.H{
				constant.HeaderMessage: constant.HeaderResponseURLDeleted,
			})
		}
	}
}

// validateAndDeleteURL validates the ID and URL and performs the deletion if they are correct.
func validateAndDeleteURL(c *gin.Context, dsClient *datastore.Client) error {
	idFromPath := c.Param(constant.HeaderID) // Extract the ID from the URL path

	// Bind the JSON payload to the DeleteURLPayload struct.
	var req DeleteURLPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		return logmonitor.NewBadRequestError(constant.HeaderResponseInvalidRequestPayload, err)
	}

	// Check if the IDs match
	if idFromPath != req.ID {
		return logmonitor.NewBadRequestError(
			constant.MisMatchBetweenPathIDandPayloadIDContextLog,
			fmt.Errorf(constant.PathIDandPayloadIDDoesnotMatchContextLog))
	}

	// Validate the URL format.
	if !isValidURL(req.URL) {
		return logmonitor.NewBadRequestError(
			constant.HeaderResponseInvalidURLFormat,
			fmt.Errorf(constant.HeaderResponseInvalidURLFormat))
	}

	// Perform the delete operation.
	return deleteURL(c, dsClient, req.ID, req.URL)
}

// deleteURL verifies the provided ID and URL against the stored URL entity, and if they match, deletes the URL entity.
func deleteURL(c *gin.Context, dsClient *datastore.Client, id string, providedURL string) error {
	currentURL, err := getCurrentURL(c, dsClient, id)
	if err != nil {
		return err // getCurrentURL will return a formatted error or datastore.ErrNotFound
	}

	if currentURL.Original != providedURL {
		return fmt.Errorf(constant.URLmismatchContextLog)
	}

	return performDelete(c, dsClient, id)
}

// getCurrentURL retrieves the current URL from the datastore and checks for errors.
func getCurrentURL(c *gin.Context, dsClient *datastore.Client, id string) (*datastore.URL, error) {
	currentURL, err := datastore.GetURL(c, dsClient, id)
	if err != nil {
		if err == datastore.ErrNotFound {
			return nil, datastore.ErrNotFound
		}
		return nil, fmt.Errorf(constant.FailedToRetriveURLContextLog+": %v", err)
	}
	return currentURL, nil
}

// performDelete deletes the URL entity from the datastore.
func performDelete(c *gin.Context, dsClient *datastore.Client, id string) error {
	if err := datastore.DeleteURL(c, dsClient, id); err != nil {
		return fmt.Errorf(constant.FailedToDeletedURLContextLog+": %v", err)
	}
	return nil
}

// saveURL saves the URL and its identifier to the datastore.
func saveURL(c *gin.Context, dsClient *datastore.Client, id string, originalURL string) error {
	url := &datastore.URL{
		Original: originalURL,
		ID:       id,
	}
	return datastore.SaveURL(c, dsClient, url)
}
