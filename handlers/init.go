package handlers

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/H0llyW00dzZ/go-urlshortner/datastore"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
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
	URL string `json:"url" binding:"required,url"`
}

// UpdateURLPayload defines the structure for the JSON payload when updating an existing URL.
// Fixed a bug potential leading to Exploit CWE-284 / IDOR in the json payloads, Now It's safe A long With ID.
type UpdateURLPayload struct {
	ID     string `json:"id" binding:"required"`
	OldURL string `json:"old_url" binding:"required,url"`
	NewURL string `json:"new_url" binding:"required,url"`
}

// DeleteURLPayload defines the structure for the JSON payload when deleting a URL.
type DeleteURLPayload struct {
	ID  string `json:"id" binding:"required"`
	URL string `json:"url" binding:"required,url"`
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
		panic(constant.InternelSecretEnvContextLog)
	}
}

// InternalOnly creates a middleware that restricts access to a route to internal services only.
// It checks for a specific header containing a secret value that should match an environment
// variable to allow the request to proceed. If the secret does not match or is not provided,
// the request is aborted with a 403 Forbidden status.
func InternalOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check the request header against the expected secret value.
		if c.GetHeader(constant.HeaderXinternalSecret) != internalSecretValue {
			// If the header does not match the expected secret, abort the request.
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				constant.HeaderResponseError: constant.HeaderResponseForbidden,
			})
			return
		}

		// If the header matches, proceed with the request.
		c.Next()
	}
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

// generateShortID generates a short identifier for the URL.
//
// The generateShortID function is responsible for generating a unique short ID.
//
// If the generated ID is not unique, it will keep trying until it finds a unique one,
// otherwise it will return an error indicate that your machine is bad.
func generateShortID(ctx context.Context, dsClient *datastore.Client) (string, error) {

	id, err := shortid.GenerateUnique(ctx, dsClient, 5)
	if err != nil {
		return "", err // If there's an error generating the ID, return it immediately.
	}
	return id, nil // If the ID is unique, return it.
}
