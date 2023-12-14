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

	cloudDatastore "cloud.google.com/go/datastore"
	localDatastore "github.com/H0llyW00dzZ/go-urlshortner/datastore"
	"github.com/H0llyW00dzZ/go-urlshortner/shortid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Logger *zap.Logger

// InternalOnly creates a middleware that restricts access to a route to internal services only.
// It checks for a specific header containing a secret value that should match an environment
// variable to allow the request to proceed. If the secret does not match or is not provided,
// the request is aborted with a 403 Forbidden status.
func InternalOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the secret value from an environment variable.
		expectedSecretValue := os.Getenv("INTERNAL_SECRET_VALUE")
		if expectedSecretValue == "" {
			// If the environment variable is not set, abort and report an internal server error.
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal configuration error"})
			return
		}

		// Check the request header against the expected secret value.
		if c.GetHeader("X-Internal-Secret") != expectedSecretValue {
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
// The base path for the handlers can be customized via the CUSTOM_BASE_PATH environment variable.
// If CUSTOM_BASE_PATH is not set, the default base path "/" is used.
func RegisterHandlersGin(router *gin.Engine, dsClient *cloudDatastore.Client) {
	// Retrieve the custom base path from an environment variable or use "/" as default.
	basePath := os.Getenv("CUSTOM_BASE_PATH")
	if basePath == "" {
		basePath = "/"
	}

	// Ensure the basePath is correctly formatted.
	if basePath[len(basePath)-1:] != "/" {
		basePath += "/"
	}

	// Register handlers with the custom or default base path.
	// For example, if CUSTOM_BASE_PATH is "/api/", the GET route will be "/api/:id" and
	// the POST route will be "/api/".
	router.GET(basePath+":id", getURLHandlerGin(dsClient))
	router.POST(basePath, InternalOnly(), postURLHandlerGin(dsClient))
}

// getURLHandlerGin returns a Gin handler function that retrieves and redirects to the original
// URL based on a short identifier provided in the request path. If the identifier is not found
// or an error occurs, the handler responds with the appropriate HTTP status code and error message.
func getURLHandlerGin(dsClient *cloudDatastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Assuming localDatastore.GetURL is a function that correctly handles datastore operations.
		url, err := localDatastore.GetURL(c, dsClient, id)
		if err != nil {
			if err == localDatastore.ErrNotFound {
				// Entity not found
				localDatastore.Logger.Warn("URL not found", zap.String("id", id))
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "URL not found"})
				return
			} else {
				// Some other error occurred
				localDatastore.Logger.Error("Failed to get URL", zap.String("id", id), zap.Error(err))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
		}

		// Check if URL is nil after the GetURL call
		if url == nil {
			localDatastore.Logger.Error("URL is nil after GetURL call", zap.String("id", id))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Redirect to the original URL if no error occurred and URL is not nil.
		c.Redirect(http.StatusFound, url.Original)
	}
}

// postURLHandlerGin returns a Gin handler function that handles the creation of a new shortened
// URL. It expects a JSON payload with the original URL, generates a short identifier, and stores
// the mapping in Google Cloud Datastore. If successful, it returns the generated identifier and
// the shortened URL; otherwise, it responds with an error.
func postURLHandlerGin(dsClient *cloudDatastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			URL string `json:"url"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			Logger.Error("Invalid request", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		id, err := shortid.Generate(5) // Generate a 5-character ID
		if err != nil {
			Logger.Error("Failed to generate ID", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		key := cloudDatastore.NameKey("urlz", id, nil)
		url := localDatastore.URL{
			Original: req.URL,
			ID:       id,
		}
		if _, err := dsClient.Put(c, key, &url); err != nil {
			Logger.Error("Failed to save URL", zap.String("id", id), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Automatically detect the base URL
		scheme := c.GetHeader("X-Forwarded-Proto") // Check for the X-Forwarded-Proto header first
		if scheme == "" {
			// Fallback to checking if the TLS is not nil
			if c.Request.TLS != nil {
				scheme = "https"
			} else {
				scheme = "http"
			}
		}
		baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)

		fullShortenedURL := baseURL + "/" + id
		c.JSON(http.StatusOK, gin.H{"id": id, "shortened_url": fullShortenedURL})
	}
}
