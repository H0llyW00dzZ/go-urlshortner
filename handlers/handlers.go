// Package handlers defines HTTP handlers and middleware for the URL shortener service.
// It provides functionality to shorten URLs and redirect to original URLs based on the
// generated short identifiers. The package uses Google Cloud Datastore for storage of
// the URL mappings and leverages middleware to restrict access to certain operations.
//
// Copyright (c) 2023 H0llyW00dzZ
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
type UpdateURLPayload struct {
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
	internalSecretValue = os.Getenv("INTERNAL_SECRET_VALUE")
	if internalSecretValue == "" {
		panic("INTERNAL_SECRET_VALUE is not set")
	}
}

// InternalOnly creates a middleware that restricts access to a route to internal services only.
// It checks for a specific header containing a secret value that should match an environment
// variable to allow the request to proceed. If the secret does not match or is not provided,
// the request is aborted with a 403 Forbidden status.
func InternalOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check the request header against the expected secret value.
		if c.GetHeader("X-Internal-Secret") != internalSecretValue {
			// If the header does not match the expected secret, abort the request.
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
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
		id := c.Param("id")

		// Assuming datastore.GetURL is a function that correctly handles datastore operations.
		url, err := datastore.GetURL(c, dsClient, id)
		if err != nil {
			logFields := []zap.Field{
				zap.String("internal", "Datastore"),
				zap.String("operation", "getURL"),
				zap.String("id", id),
				zap.Error(err),
			}
			if err == datastore.ErrNotFound {
				// Entity not found
				// Friendly logger not found message for the devops that always monitoring logs
				Logger.Info("URL not found", logFields...)
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "URL not found"})
				return
			} else {
				Logger.Error("Failed to get URL", logFields...)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
		}

		// Check if URL is nil after the GetURL call
		if url == nil {
			Logger.Error("URL is nil after GetURL call",
				zap.String("operation", "getURL"),
				zap.String("id", id),
			)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

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
			handleError(c, "Invalid request payload", http.StatusBadRequest, err)
			return
		}

		// Generate a short identifier for the URL.
		id, err := generateShortID()
		if err != nil {
			handleError(c, "Failed to generate ID", http.StatusInternalServerError, err)
			return
		}

		// Save the URL with the generated identifier into the datastore.
		if err := saveURL(c, dsClient, id, url); err != nil {
			handleError(c, "Failed to save URL", http.StatusInternalServerError, err)
			return
		}

		// Construct the full shortened URL and return it in the response.
		fullShortenedURL := constructFullShortenedURL(c, id)
		c.JSON(http.StatusOK, gin.H{"id": id, "shortened_url": fullShortenedURL})
	}
}

// editURLHandlerGin returns a Gin handler function that handles the updating of an existing shortened URL.
func editURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the ID from the URL path parameter.
		id := c.Param("id")

		// Bind the JSON payload to the UpdateURLPayload struct.
		req, err := bindUpdatePayload(c)
		if err != nil {
			handleError(c, "Invalid request", http.StatusBadRequest, err)
			return
		}

		// Perform the update operation.
		if err := updateURL(c, dsClient, id, req); err != nil {
			handleError(c, err.Error(), http.StatusInternalServerError, err)
			return
		}

		// Respond with the updated URL information.
		respondWithUpdatedURL(c, id)
	}
}

// bindUpdatePayload binds the JSON payload to the UpdateURLPayload struct and validates the new URL format.
func bindUpdatePayload(c *gin.Context) (UpdateURLPayload, error) {
	var req UpdateURLPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, err
	}

	if req.NewURL == "" || !isValidURL(req.NewURL) {
		return req, fmt.Errorf("invalid new URL format")
	}

	return req, nil
}

// updateURL retrieves the current URL, verifies it against the provided old URL, and updates it with the new URL.
// It returns an error with a message suitable for HTTP response if any step fails.
func updateURL(c *gin.Context, dsClient *datastore.Client, id string, req UpdateURLPayload) error {
	// Retrieve the current URL to ensure it matches the provided old URL.
	currentURL, err := datastore.GetURL(c, dsClient, id)
	if err != nil {
		return fmt.Errorf("failed to retrieve URL")
	}
	if currentURL.Original != req.OldURL {
		return fmt.Errorf("URL mismatch")
	}

	// Update the URL in the datastore with the new URL.
	if err := datastore.UpdateURL(c, dsClient, id, req.NewURL); err != nil {
		return fmt.Errorf("failed to update URL")
	}

	return nil
}

// respondWithUpdatedURL constructs and sends a JSON response with the updated URL information.
func respondWithUpdatedURL(c *gin.Context, id string) {
	fullShortenedURL := constructFullShortenedURL(c, id)
	c.JSON(http.StatusOK, gin.H{
		"id":            id,
		"shortened_url": fullShortenedURL,
		"status":        "URL updated successfully",
	})
}

// extractURL extracts the original URL from the JSON payload in the request.
func extractURL(c *gin.Context) (string, error) {
	var req CreateURLPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		Logger.Error("Invalid request - JSON binding error", zap.Error(err))
		return "", err
	}

	// Check if the URL is in a valid format.
	if req.URL == "" || !isValidURL(req.URL) {
		Logger.Error("Invalid URL format", zap.String("url", req.URL))
		return "", fmt.Errorf("invalid URL format")
	}

	return req.URL, nil
}

// deleteURLHandlerGin returns a Gin handler function that handles the deletion of an existing shortened URL.
func deleteURLHandlerGin(dsClient *datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validateAndDeleteURL(c, dsClient); err != nil {
			handleDeletionError(c, err)
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "URL deleted successfully"})
		}
	}
}

// handleDeletionError handles errors that occur during the URL deletion process.
func handleDeletionError(c *gin.Context, err error) {
	if badRequestErr, ok := err.(*logmonitor.BadRequestError); ok {
		handleError(c, badRequestErr.UserMessage, http.StatusBadRequest, badRequestErr.Err)
	} else if err == datastore.ErrNotFound {
		Logger.Info("URL not found for deletion", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
	} else {
		Logger.Error("Failed to delete URL", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}

// validateAndDeleteURL validates the ID and URL and performs the deletion if they are correct.
func validateAndDeleteURL(c *gin.Context, dsClient *datastore.Client) error {
	idFromPath := c.Param("id") // Extract the ID from the URL path

	// Bind the JSON payload to the DeleteURLPayload struct.
	var req DeleteURLPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		return logmonitor.NewBadRequestError("Invalid request payload", err)
	}

	// Check if the IDs match
	if idFromPath != req.ID {
		return logmonitor.NewBadRequestError("Mismatch between path ID and payload ID", fmt.Errorf("path ID and payload ID do not match"))
	}

	// Validate the URL format.
	if !isValidURL(req.URL) {
		return logmonitor.NewBadRequestError("Invalid URL format", fmt.Errorf("invalid URL format"))
	}

	// Perform the delete operation.
	return deleteURL(c, dsClient, req.ID, req.URL)
}

// deleteURL verifies the provided ID and URL against the stored URL entity, and if they match, deletes the URL entity.
func deleteURL(c *gin.Context, dsClient *datastore.Client, id string, providedURL string) error {
	// Retrieve the current URL to ensure it matches the provided URL.
	currentURL, err := datastore.GetURL(c, dsClient, id)
	if err != nil {
		return fmt.Errorf("failed to retrieve URL: %v", err)
	}
	if currentURL.Original != providedURL {
		return fmt.Errorf("URL mismatch")
	}

	// Delete the URL in the datastore.
	if err := datastore.DeleteURL(c, dsClient, id); err != nil {
		return fmt.Errorf("failed to delete URL: %v", err)
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
	scheme := c.GetHeader("X-Forwarded-Proto")
	if scheme == "" {
		// Fallback to checking the TLS property of the request if the header is not set.
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
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
	// Use different log levels based on the status code
	switch {
	case statusCode >= 500: // 5xx errors are still logged as errors
		logmonitor.Logger.Error(message, zap.Error(err))
	case statusCode >= 400: // 4xx errors are logged as warnings
		logmonitor.Logger.Info(message, zap.Error(err))
	default: // All other status codes are logged as info
		logmonitor.Logger.Info(message, zap.Error(err))
	}
	c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
}
