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
	"os"
	"strings"

	"github.com/H0llyW00dzZ/go-urlshortner/datastore"
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

func init() {
	config := zap.NewDevelopmentConfig()
	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}
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

// RegisterHandlersGin registers the HTTP handlers for the URL shortener service using the Gin
// web framework. It sets up the routes for retrieving and creating shortened URLs and applies
// the InternalOnly middleware to the POST route to protect it from public access.
func RegisterHandlersGin(router *gin.Engine, datastoreClient *datastore.Client) {
	// Register handlers with the custom or default base path.
	// For example, if CUSTOM_BASE_PATH is "/api/", the GET route will be "/api/:id" and
	// the POST route will be "/api/".
	router.GET(basePath+":id", getURLHandlerGin(datastoreClient))
	router.POST(basePath, InternalOnly(), postURLHandlerGin(datastoreClient))
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
				zap.String("operation", "getURL"),
				zap.String("id", id),
				zap.Error(err),
			}
			if err == datastore.ErrNotFound {
				// Entity not found
				Logger.Warn("URL not found", logFields...)
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
		// Extract the original URL from the request body.
		url, err := extractURL(c)
		if err != nil {
			handleError(c, "Invalid request", http.StatusBadRequest, err)
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

// extractURL extracts the original URL from the JSON payload in the request.
func extractURL(c *gin.Context) (string, error) {
	var req struct {
		URL string `json:"url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Logger.Error("Invalid request", zap.Error(err))
		return "", err
	}
	return req.URL, nil
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
	Logger.Error(message, zap.Error(err))
	c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
}
