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
)

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
func RegisterHandlersGin(router *gin.Engine, dsClient *cloudDatastore.Client) {
	router.GET("/:id", getURLHandlerGin(dsClient))

	// Apply the InternalOnly middleware to the POST route.
	router.POST("/", InternalOnly(), postURLHandlerGin(dsClient))
}

// getURLHandlerGin returns a Gin handler function that retrieves and redirects to the original
// URL based on a short identifier provided in the request path. If the identifier is not found
// or an error occurs, the handler responds with the appropriate HTTP status code and error message.
func getURLHandlerGin(dsClient *cloudDatastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		key := cloudDatastore.NameKey("urlz", id, nil)
		var url localDatastore.URL
		if err := dsClient.Get(c, key, &url); err != nil {
			fmt.Printf("Failed to get URL: %v\n", err)
			if err == cloudDatastore.ErrNoSuchEntity {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "URL not found"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
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
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		id, err := shortid.Generate(5) // Generate a 5-character ID
		if err != nil {
			fmt.Printf("Failed to generate ID: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		key := cloudDatastore.NameKey("urlz", id, nil)
		url := localDatastore.URL{
			Original: req.URL,
			ID:       id,
		}
		if _, err := dsClient.Put(c, key, &url); err != nil {
			fmt.Printf("Failed to save URL: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Automatically detect the base URL
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)

		fullShortenedURL := baseURL + "/" + id
		c.JSON(http.StatusOK, gin.H{"id": id, "shortened_url": fullShortenedURL})
	}
}
