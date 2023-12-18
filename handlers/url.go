package handlers

import (
	"fmt"
	"net/http"

	"github.com/H0llyW00dzZ/go-urlshortner/datastore"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"github.com/gin-gonic/gin"
)

// getURLHandlerGin returns a Gin handler function that retrieves and redirects to the original
// URL based on a short identifier provided in the request path. If the identifier is not found
// or an error occurs, the handler responds with the appropriate HTTP status code and error message.
func getURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the identifier from the request path.
		id := c.Param(constant.HeaderID)

		// Attempt to retrieve the URL from the datastore using the provided identifier.
		url, err := datastore.GetURL(c, dsClient, id)

		// Error handling block.
		if err != nil {
			// If the URL is not found, log the event and return a 404 Not Found response.
			if err == datastore.ErrNotFound {
				// Use the centralized logging function to log the 'URL not found' event.
				LogURLNotFound(id, err)
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					constant.HeaderResponseError: constant.URLnotfoundContextLog,
				})
			} else {
				// For any other errors, log the internal error event and return a 500 Internal Server Error response.
				LogInternalError("getURL", id, err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					constant.HeaderResponseError: constant.HeaderResponseInternalServerError,
				})
			}
			return
		}

		// If the retrieved URL is nil, log the error and return a 500 Internal Server Error response.
		if url == nil {
			LogInternalError("getURL", id, fmt.Errorf(constant.URLisNilContextLog))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				constant.HeaderResponseError: constant.HeaderResponseInternalServerError,
			})
			return
		}

		// Log the successful retrieval of the URL.
		LogURLRetrievalSuccess(id)

		// Redirect the client to the original URL associated with the identifier.
		c.Redirect(http.StatusFound, url.Original)
	}
}

// postURLHandlerGin returns a Gin handler function that handles the creation of a new shortened
// URL. It expects a JSON payload with the original URL, generates a short identifier, and stores
// the mapping in the datastore. If successful, it returns the generated identifier and
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
		id, err := generateShortID(c.Request.Context(), dsClient)
		if err != nil {
			handleError(c, constant.HeaderResponseFailedtoGenerateID, http.StatusInternalServerError, err)
			return
		}

		// Save the URL with the generated identifier into the datastore.
		if err := saveURL(c, dsClient, id, url); err != nil {
			handleError(c, constant.HeaderResponseFailedtoSaveURL, http.StatusInternalServerError, err)
			return
		}

		// Use the centralized logging function to log the successful shortening of the URL.
		LogURLShortened(id)

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
		// Use the centralized logging function to log the bad request error.
		LogBadRequestError("validateUpdateRequest", err)
		return "", req, err
	}

	// Additional validation for the ID in the URL and the ID in the payload
	if pathID != req.ID {
		err := fmt.Errorf(constant.PathIDandPayloadIDDoesnotMatchContextLog)
		// Use the centralized logging function to log the mismatch error.
		LogMismatchError(pathID)
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
		LogMismatchError(id)
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
		// Replace the direct logger call with a centralized logging function
		LogBadRequestError("extractURL", err)
		return "", err
	}

	// Check if the URL is in a valid format.
	if req.URL == "" || !isValidURL(req.URL) {
		// Replace the direct logger call with a centralized logging function
		LogInvalidURLFormat(req.URL)
		return "", fmt.Errorf(constant.HeaderResponseInvalidURLFormat)
	}

	return req.URL, nil
}

// deleteURLHandlerGin returns a Gin handler function that handles the deletion of an existing shortened URL.
func deleteURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param(constant.HeaderID)
		if err := validateAndDeleteURL(c, dsClient); err != nil {
			// Use the centralized logging function to log the deletion error.
			LogDeletionError(id, err)
			handleDeletionError(c, err)
		} else {
			// Use the centralized logging function to log the successful deletion.
			LogURLDeletionSuccess(id)
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
		LogBadRequestError("deleteURL", err) // Log the bad request error
		return fmt.Errorf(constant.HeaderResponseInvalidRequestPayload+": %v", err)
	}

	// Check if the IDs match
	if idFromPath != req.ID {
		LogMismatchError(idFromPath) // Log the mismatch error
		return fmt.Errorf(constant.PathIDandPayloadIDDoesnotMatchContextLog)
	}

	// Validate the URL format.
	if !isValidURL(req.URL) {
		LogInvalidURLFormat(req.URL) // Log the invalid URL format error
		return fmt.Errorf(constant.HeaderResponseInvalidURLFormat)
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
