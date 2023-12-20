package handlers

import (
	"fmt"
	"net/http"

	"github.com/H0llyW00dzZ/go-urlshortner/datastore"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// getURLHandlerGin returns a Gin handler function that retrieves and redirects to the original
// URL based on a short identifier provided in the request path. If the identifier is not found
// or an error occurs, the handler responds with the appropriate HTTP status code and error message.
func getURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Apply the rate limiter first.
		if !applyRateLimit(c) {
			// If the rate limit is exceeded, we should not continue processing the request.
			logAttemptToRetrieve(id)
			return
		}

		id := c.Param(constant.HeaderID)
		url, err := datastore.GetURL(c, dsClient, id)
		if err != nil {
			handleGetURLError(c, id, err)
			return
		}

		if url == nil {
			LogInternalError(operation_getURL, id, fmt.Errorf(constant.URLisNilContextLog))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				constant.HeaderResponseError: constant.HeaderResponseInternalServerError,
			})
			return
		}

		LogURLRetrievalSuccess(id)
		c.Redirect(http.StatusFound, url.Original)
	}
}

// handleGetURLError centralizes the error handling for the getURLHandlerGin function.
func handleGetURLError(c *gin.Context, id string, err error) {
	if err == datastore.ErrNotFound {
		LogURLNotFound(id, err)
		// Respond with 404 Not Found, as this is the correct status for a missing resource.
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			constant.HeaderResponseError: constant.URLnotfoundContextLog,
		})
		// For any other errors, log the internal error event and return a 500 Internal Server Error response.
	} else {
		LogInternalError(operation_getURL, id, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			constant.HeaderResponseError: constant.HeaderResponseInternalServerError,
		})
	}
}

// applyRateLimit checks if the rate limit has been exceeded for a given identifier.
// It writes the appropriate response if the rate limit is exceeded.
//
// Note: Better keep set it as response 404 not found for this limiter since it's public access,
// for indicate that client/user are bad requesting, not the server.
func applyRateLimit(c *gin.Context) bool {
	key := c.ClientIP()
	limiter := NewRateLimiter(key, rate.Limit(5), 10)
	if !limiter.Allow() {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			constant.HeaderResponseError: constant.HeaderResponseIDandURLNotFound,
		})
		return false
	}
	return true
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
		//LogBadRequestError("validateUpdateRequest", err)
		SynclogError(c, operation_validateUpdateRequest, err) // Replaced with centralized logging function
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
		// Instead of handling the error here, we return it to the caller to handle.
		return handleRetrievalError(err, id)
	}

	if currentURL.Original != req.OldURL {
		// Return a URLMismatchError which can be handled specifically by the caller.
		return &URLMismatchError{Message: constant.URLmismatchContextLog}
	}

	logAttemptToUpdate(id)

	// Update the URL in the datastore with the new URL.
	if err := datastore.UpdateURL(c, dsClient, id, req.NewURL); err != nil {
		// Return the error to the caller to handle.
		return err
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
		// LogBadRequestError("extractURL", err)
		SynclogError(c, operation_extractURL, err) // Replaced with centralized logging function
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
		SynclogError(c, idFromPath, err) // Log the bad request error
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
	// Retrieve the current URL from the datastore.
	currentURL, err := getCurrentURL(c, dsClient, id)
	if err != nil {
		// If an error occurs, return it. getCurrentURL will return a formatted error or datastore.ErrNotFound.
		return err
	}

	// Check if the current URL matches the provided URL.
	if currentURL.Original != providedURL {
		// If they do not match, return a custom URLMismatchError instead of a generic error (known as default standart library error/fmt error),
		// which is bad for host machine and datastore when using generic error, it literally break the machine (can't imagine if there is no recovery mode lol).
		return &URLMismatchError{Message: constant.URLmismatchContextLog}
	}

	// If the URLs match, perform the deletion operation.
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
