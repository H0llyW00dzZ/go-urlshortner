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

// InternalOnly is a middleware that ensures only internal services can access the route.
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

func RegisterHandlersGin(router *gin.Engine, dsClient *cloudDatastore.Client) {
	router.GET("/:id", getURLHandlerGin(dsClient))

	// Apply the InternalOnly middleware to the POST route.
	router.POST("/", InternalOnly(), postURLHandlerGin(dsClient))
}

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
		c.JSON(http.StatusOK, gin.H{"id": id})
	}
}
