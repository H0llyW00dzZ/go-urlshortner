package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/H0llyW00dzZ/go-urlshortner/datastore"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor"
	"github.com/H0llyW00dzZ/go-urlshortner/shortid"
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

// CreateURLPayload defines the structure for the JSON payload when creating a new URL.
// It contains a single field, URL, which is the original URL to be shortened.
type CreateURLPayload struct {
	URL string `json:"url"`
}

// UpdateURLPayload defines the structure for the JSON payload when updating an existing URL.
// Fixed a bug potential leading to Exploit CWE-284 / IDOR in the json payloads, Now It's safe A long With ID.
type UpdateURLPayload struct {
	ID     string `json:"id"`
	OldURL string `json:"old_url"`
	NewURL string `json:"new_url"`
}

// DeleteURLPayload defines the structure for the JSON payload when deleting a URL.
type DeleteURLPayload struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func init() {

	// Initialize the base path from an environment variable or use "/" as default.
	basePath = os.Getenv("CUSTOM_BASE_PATH")
	if basePath == "" {
		basePath = "/"
	}
	// Ensure the basePath is correctly formatted.
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	// Initialize the internal secret value from an environment variable.
	// Note: This is important and secure because it resides deep within the binary internals and should not be left unset in production.
	internalSecretValue = os.Getenv("INTERNAL_SECRET_VALUE")
	if internalSecretValue == "" {
		panic(logmonitor.InternelSecretEnvContextLog)
	}
}

// InternalOnly creates a middleware that restricts access to a route to internal services only.
// It checks for a specific header containing a secret value that should match an environment
// variable to allow the request to proceed. If the secret does not match or is not provided,
// the request is aborted with a 403 Forbidden status.
func InternalOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check the request header against the expected secret value.
		if c.GetHeader(logmonitor.HeaderXinternalSecret) != internalSecretValue {
			// If the header does not match the expected secret, abort the request.
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				logmonitor.HeaderResponseError: logmonitor.HeaderResponseForbidden,
			})
			return
		}

		// If the header matches, proceed with the request.
		c.Next()
	}
}

// isValidURL checks if the URL is in a valid format.
func isValidURL(urlStr string) bool {
	u, err := url.ParseRequestURI(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// RegisterHandlersGin registers the HTTP handlers for the URL shortener service using the Gin
// web framework. It sets up the routes for retrieving, creating, and updating shortened URLs.
// The InternalOnly middleware is applied to the POST and PUT routes to protect them from public access.
func RegisterHandlersGin(router *gin.Engine, datastoreClient *datastore.Client) {
	// Register handlers with the custom or default base path.
	// For example, if CUSTOM_BASE_PATH is "/api/", the GET route will be "/api/:id",
	// the POST route will be "/api/", and the PUT route will be "/api/:id".
	router.GET(basePath+":id", getURLHandlerGin(datastoreClient))
	router.POST(basePath, InternalOnly(), postURLHandlerGin(datastoreClient))
	router.PUT(basePath+":id", InternalOnly(), editURLHandlerGin(datastoreClient))      // New PUT route for editing URLs
	router.DELETE(basePath+":id", InternalOnly(), deleteURLHandlerGin(datastoreClient)) // New DELETE route for deleting URLs
}

// getURLHandlerGin returns a Gin handler function that retrieves and redirects to the original
// URL based on a short identifier provided in the request path. If the identifier is not found
// or an error occurs, the handler responds with the appropriate HTTP status code and error message.
func getURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param(logmonitor.HeaderID)

		// Assuming datastore.GetURL is a function that correctly handles datastore operations.
		url, err := datastore.GetURL(c, dsClient, id)
		// Declare logFields here so it's accessible throughout the function scope
		logFields := logmonitor.CreateLogFields("getURL",
			logmonitor.WithComponent(logmonitor.ComponentNoSQL), // Use the constant for the component
			logmonitor.WithID(id),
			logmonitor.WithError(err), // Include the error here, but it will be nil if there's no error
		)

		if err != nil {
			if err == datastore.ErrNotFound {
				logmonitor.Logger.Info(logmonitor.GetBackEmoji+"  "+logmonitor.UrlshortenerEmoji+"  "+logmonitor.URLnotfoundContextLog, logFields...)
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					logmonitor.HeaderResponseError: logmonitor.URLnotfoundContextLog,
				})
			} else {
				logmonitor.Logger.Error(logmonitor.SosEmoji+"  "+logmonitor.WarningEmoji+"  "+logmonitor.FailedToGetURLContextLog, logFields...)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					logmonitor.HeaderResponseError: logmonitor.HeaderResponseInternalServerError,
				})
			}
			return
		}

		// Check if URL is nil after the GetURL call
		if url == nil {
			// Use the logmonitor's logging function for consistency
			logmonitor.Logger.Error(logmonitor.SosEmoji+"  "+logmonitor.WarningEmoji+"  "+logmonitor.URLisNilContextLog, logFields...)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				logmonitor.HeaderResponseError: logmonitor.HeaderResponseInternalServerError,
			})
			return
		}

		// If there's no error and you're logging a successful retrieval, use the same logFields
		logmonitor.Logger.Info(logmonitor.UrlshortenerEmoji+"  "+logmonitor.RedirectEmoji+"  "+logmonitor.SuccessEmoji+"  "+logmonitor.URLRetriveContextLog, logFields...)
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
			handleError(c, logmonitor.HeaderResponseInvalidRequestPayload, http.StatusBadRequest, err)
			return
		}

		// Generate a short identifier for the URL.
		id, err := generateShortID()
		if err != nil {
			handleError(c, logmonitor.HeaderResponseFailedtoGenerateID, http.StatusInternalServerError, err)
			return
		}

		// Save the URL with the generated identifier into the datastore.
		if err := saveURL(c, dsClient, id, url); err != nil {
			handleError(c, logmonitor.HeaderResponseFailedtoSaveURL, http.StatusInternalServerError, err)
			return
		}

		logFields := logmonitor.CreateLogFields("postURL",
			logmonitor.WithComponent(logmonitor.ComponentNoSQL), // Use the constant for the component
			logmonitor.WithID(id),
		)

		logmonitor.Logger.Info(logmonitor.UrlshortenerEmoji+"  "+logmonitor.SuccessEmoji+"  "+logmonitor.URLShorteneredContextLog, logFields...)

		// Construct the full shortened URL and return it in the response.
		fullShortenedURL := constructFullShortenedURL(c, id)
		c.JSON(http.StatusOK, gin.H{
			logmonitor.HeaderID: id, logmonitor.HeaderResponseshortened_url: fullShortenedURL,
		})
	}
}

// editURLHandlerGin returns a Gin handler function that handles the updating of an existing shortened URL.
func editURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		pathID, req, err := validateUpdateRequest(c)
		if err != nil {
			handleError(c, err.Error(), http.StatusBadRequest, err)
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
	pathID = c.Param(logmonitor.HeaderID)

	if err := c.ShouldBindJSON(&req); err != nil {
		return "", req, err
	}

	if pathID != req.ID {
		return "", req, fmt.Errorf(logmonitor.MisMatchBetweenPathIDandPayloadIDContextLog)
	}

	return pathID, req, nil
}

// handleUpdateError handles errors that occur during the URL update process.
func handleUpdateError(c *gin.Context, id string, err error) {
	logFields := logmonitor.CreateLogFields("editURL",
		logmonitor.WithComponent(logmonitor.ComponentNoSQL),
		logmonitor.WithID(id),
		logmonitor.WithError(err),
	)

	if strings.Contains(err.Error(), logmonitor.URLmismatchContextLog) {
		logmonitor.Logger.Info(logmonitor.UrlshortenerEmoji+"  "+logmonitor.UpdateEmoji+"  "+logmonitor.ErrorEmoji+"  "+logmonitor.URLmismatchContextLog, logFields...)
		c.JSON(http.StatusBadRequest, gin.H{
			logmonitor.HeaderResponseError: err.Error(),
		})
		return
	}

	logmonitor.Logger.Info(logmonitor.AlertEmoji+"  "+logmonitor.WarningEmoji+"  "+logmonitor.FailedToUpdateURLContextLog, logFields...)
	c.JSON(http.StatusInternalServerError, gin.H{
		logmonitor.HeaderResponseError: err.Error(),
	})
}

// bindUpdatePayload binds the JSON payload to the UpdateURLPayload struct and validates the new URL format.
func bindUpdatePayload(c *gin.Context) (UpdateURLPayload, error) {
	var req UpdateURLPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, err
	}

	if req.NewURL == "" || !isValidURL(req.NewURL) {
		return req, fmt.Errorf(logmonitor.InvalidNewURLFormatContextLog)
	}

	return req, nil
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
		return fmt.Errorf(logmonitor.URLmismatchContextLog)
	}

	logAttemptToUpdate(id)

	// Update the URL in the datastore with the new URL.
	if err := datastore.UpdateURL(c, dsClient, id, req.NewURL); err != nil {
		handleUpdateError(c, id, err) // handle the error, but don't expect a return value
		return fmt.Errorf(logmonitor.FailedToUpdateURLContextLog+" %v", err)
	}

	logSuccessfulUpdate(id)

	return nil
}

// logAttemptToRetrieve logs an informational message indicating an attempt to retrieve the current URL by ID.
func logAttemptToRetrieve(id string) {
	logFields := createLogFields(id)
	logmonitor.Logger.Info(logmonitor.AlertEmoji+"  "+logmonitor.WarningEmoji+"  "+logmonitor.InfoAttemptingToRetrieveTheCurrentURL, logFields...)
}

// handleRetrievalError logs an error message for a failed retrieval attempt and returns a formatted error.
// If the error is a 'not found' error, it logs a specific message for that case.
func handleRetrievalError(err error, id string) error {
	logFields := createLogFields(id)
	if err == datastore.ErrNotFound {
		logmonitor.Logger.Info(logmonitor.GetBackEmoji+"  "+logmonitor.UrlshortenerEmoji+"  "+logmonitor.URLnotfoundContextLog, logFields...)
		return fmt.Errorf(logmonitor.URLnotfoundContextLog)
	}
	logmonitor.Logger.Error(logmonitor.FailedToRetriveURLContextLog+": "+err.Error(), logFields...)
	return fmt.Errorf(logmonitor.FailedToRetriveURLContextLog+": %v", err)
}

// logMismatchError logs an informational message indicating a mismatch error during URL update process.
func logMismatchError(id string) {
	logFields := createLogFields(id)
	logmonitor.Logger.Info(logmonitor.UrlshortenerEmoji+"  "+logmonitor.ErrorEmoji+"  "+logmonitor.URLmismatchContextLog, logFields...)
}

// logAttemptToUpdate logs an informational message indicating an attempt to update a URL in the datastore.
func logAttemptToUpdate(id string) {
	logFields := createLogFields(id)
	logmonitor.Logger.Info(logmonitor.AlertEmoji+"  "+logmonitor.WarningEmoji+"  "+datastore.InfoAttemptingToUpdateURLInDatastore, logFields...)
}

// logSuccessfulUpdate logs an informational message indicating a successful update of a URL in the datastore.
func logSuccessfulUpdate(id string) {
	logFields := createLogFields(id)
	logmonitor.Logger.Info(logmonitor.UrlshortenerEmoji+"  "+logmonitor.UpdateEmoji+"  "+logmonitor.SuccessEmoji+"  "+datastore.InfoUpdateSuccessful, logFields...)
}

// createLogFields generates a slice of zap.Field containing common log fields for the updateURL operation.
func createLogFields(id string) []zap.Field {
	return logmonitor.CreateLogFields("updateURL",
		logmonitor.WithComponent(logmonitor.ComponentNoSQL),
		logmonitor.WithID(id),
	)
}

// respondWithUpdatedURL constructs and sends a JSON response with the updated URL information.
func respondWithUpdatedURL(c *gin.Context, id string) {
	fullShortenedURL := constructFullShortenedURL(c, id)
	c.JSON(http.StatusOK, gin.H{
		logmonitor.HeaderID:                    id,
		logmonitor.HeaderResponseshortened_url: fullShortenedURL,
		logmonitor.HeaderResponseStatus:        logmonitor.HeaderResponseURlUpdated,
	})
}

// extractURL extracts the original URL from the JSON payload in the request.
func extractURL(c *gin.Context) (string, error) {
	var req CreateURLPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		Logger.Info(logmonitor.AlertEmoji+"  "+logmonitor.WarningEmoji+"  "+logmonitor.HeaderResponseInvalidRequestJSONBinding, zap.Error(err))
		return "", err
	}

	// Check if the URL is in a valid format.
	if req.URL == "" || !isValidURL(req.URL) {
		Logger.Info(logmonitor.AlertEmoji+"  "+logmonitor.WarningEmoji+"  "+logmonitor.HeaderResponseInvalidURLFormat, zap.String("url", req.URL))
		return "", fmt.Errorf(logmonitor.HeaderResponseInvalidURLFormat)
	}

	return req.URL, nil
}

// deleteURLHandlerGin returns a Gin handler function that handles the deletion of an existing shortened URL.
func deleteURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use the centralized logging function from logmonitor package
		logFields := logmonitor.CreateLogFields("deleteURL",
			logmonitor.WithComponent(logmonitor.ComponentNoSQL), // Use the constant for the component
			logmonitor.WithID(c.Param(logmonitor.HeaderID)),
		)
		if err := validateAndDeleteURL(c, dsClient); err != nil {
			handleDeletionError(c, err)
		} else {
			logmonitor.Logger.Info(logmonitor.DeleteEmoji+"  "+logmonitor.UrlshortenerEmoji+"  "+logmonitor.SuccessEmoji+"  "+logmonitor.HeaderResponseURLDeleted, logFields...)
			c.JSON(http.StatusOK, gin.H{
				logmonitor.HeaderMessage: logmonitor.HeaderResponseURLDeleted,
			})
		}
	}
}

// handleDeletionError handles errors that occur during the URL deletion process.
func handleDeletionError(c *gin.Context, err error) {
	id := c.Param(logmonitor.HeaderID)
	// Use the centralized logging function from logmonitor package
	logFields := logmonitor.CreateLogFields("deleteURL",
		logmonitor.WithComponent(logmonitor.ComponentNoSQL), // Use the constant for the component
		logmonitor.WithID(id),
		logmonitor.WithError(err),
	)
	// Fix internal issue now it's stable
	switch {
	case err == datastore.ErrNotFound:
		logmonitor.Logger.Info(logmonitor.AlertEmoji+"  "+logmonitor.WarningEmoji+"  "+logmonitor.NoURLIDContextLog, logFields...)
		c.JSON(http.StatusNotFound, gin.H{
			logmonitor.HeaderResponseError: logmonitor.HeaderResponseIDandURLNotFound,
		})
	case strings.Contains(err.Error(), logmonitor.URLmismatchContextLog):
		logmonitor.Logger.Info(logmonitor.AlertEmoji+"  "+logmonitor.WarningEmoji+"  "+logmonitor.URLmismatchContextLog, logFields...)
		c.JSON(http.StatusBadRequest, gin.H{
			logmonitor.HeaderResponseError: logmonitor.URLmismatchContextLog,
		})
	case err.(*logmonitor.BadRequestError) != nil:
		badRequestErr := err.(*logmonitor.BadRequestError)
		logmonitor.Logger.Info(logmonitor.AlertEmoji+"  "+logmonitor.WarningEmoji+"  "+logmonitor.FailedToValidateURLContextLog, logFields...)
		c.JSON(http.StatusBadRequest, gin.H{
			logmonitor.HeaderResponseError: badRequestErr.UserMessage,
		})
	default:
		logmonitor.Logger.Error(logmonitor.SosEmoji+"  "+logmonitor.WarningEmoji+"  "+logmonitor.FailedToDeletedURLContextLog, logFields...)
		c.JSON(http.StatusInternalServerError, gin.H{
			logmonitor.HeaderResponseError: logmonitor.HeaderResponseInternalServerError,
		})
	}
}

// validateAndDeleteURL validates the ID and URL and performs the deletion if they are correct.
func validateAndDeleteURL(c *gin.Context, dsClient *datastore.Client) error {
	idFromPath := c.Param(logmonitor.HeaderID) // Extract the ID from the URL path

	// Bind the JSON payload to the DeleteURLPayload struct.
	var req DeleteURLPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		return logmonitor.NewBadRequestError(logmonitor.HeaderResponseInvalidRequestPayload, err)
	}

	// Check if the IDs match
	if idFromPath != req.ID {
		return logmonitor.NewBadRequestError(
			logmonitor.MisMatchBetweenPathIDandPayloadIDContextLog,
			fmt.Errorf(logmonitor.PathIDandPayloadIDDoesnotMatchContextLog))
	}

	// Validate the URL format.
	if !isValidURL(req.URL) {
		return logmonitor.NewBadRequestError(
			logmonitor.HeaderResponseInvalidURLFormat,
			fmt.Errorf(logmonitor.HeaderResponseInvalidURLFormat))
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
		return fmt.Errorf(logmonitor.URLmismatchContextLog)
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
		return nil, fmt.Errorf(logmonitor.FailedToRetriveURLContextLog+": %v", err)
	}
	return currentURL, nil
}

// performDelete deletes the URL entity from the datastore.
func performDelete(c *gin.Context, dsClient *datastore.Client, id string) error {
	if err := datastore.DeleteURL(c, dsClient, id); err != nil {
		return fmt.Errorf(logmonitor.FailedToDeletedURLContextLog+": %v", err)
	}
	return nil
}

// generateShortID generates a short identifier for the URL.
func generateShortID() (string, error) {
	return shortid.Generate(5)
}

// saveURL saves the URL and its identifier to the datastore.
func saveURL(c *gin.Context, dsClient *datastore.Client, id string, originalURL string) error {
	url := &datastore.URL{
		Original: originalURL,
		ID:       id,
	}
	return datastore.SaveURL(c, dsClient, url)
}

// constructFullShortenedURL constructs the full shortened URL from the request and the base path.
func constructFullShortenedURL(c *gin.Context, id string) string {
	// Check for the X-Forwarded-Proto header to determine the scheme.
	scheme := c.GetHeader(logmonitor.HeaderXProto)
	if scheme == "" {
		// Fallback to checking the TLS property of the request if the header is not set.
		if c.Request.TLS != nil {
			scheme = logmonitor.HeaderSchemeHTTPS
		} else {
			scheme = logmonitor.HeaderSchemeHTTP
		}
	}

	baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)

	// Normalize the basePath by trimming leading and trailing slashes
	normalizedBasePath := strings.Trim(basePath, "/")

	// Construct the final URL ensuring there's exactly one slash between each part
	var fullPath string
	if normalizedBasePath == "" {
		fullPath = fmt.Sprintf("%s/%s", baseURL, id)
	} else {
		fullPath = fmt.Sprintf("%s/%s/%s", baseURL, normalizedBasePath, id)
	}

	return fullPath
}

// handleError logs the error and sends a JSON response with the error message and status code.
func handleError(c *gin.Context, message string, statusCode int, err error) {
	var emoji string

	// Use different emojis based on the status code
	switch {
	case statusCode >= 500: // 5xx errors are still logged as errors
		emoji = logmonitor.ErrorEmoji
		Logger.Error(emoji+"  "+message, zap.Error(err))
	}

	c.AbortWithStatusJSON(statusCode, gin.H{
		logmonitor.HeaderResponseError: message,
	})
}
