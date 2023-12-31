package handlers

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/H0llyW00dzZ/go-urlshortner/datastore"
	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"github.com/H0llyW00dzZ/go-urlshortner/shortid"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// id is a package-level variable to store the ID of the URL.
var id string

// basePath is a package-level variable to store the base path for the handlers.
// It is set once during package initialization.
var basePath string

// internalSecretValue is a package-level variable that stores the secret value required by the InternalOnly middleware.
// It is set once during package initialization.
var internalSecretValue string

// RateLimiterStore stores the rate limiters for each client, identified by a key such as an IP address.
//
// Note: This var are contains filtered in docs indicates that explicit unreadable for human 🏴‍☠️
var RateLimiterStore = struct {
	sync.RWMutex
	limiters map[string]*rate.Limiter
}{limiters: make(map[string]*rate.Limiter)}

func init() {

	// Initialize the base path from an environment variable or use "/" as default.
	basePath = os.Getenv(CUSTOM_BASE_PATH)
	if basePath == "" {
		basePath = PathObjectBasePath
	}
	// Ensure the basePath is correctly formatted.
	if !strings.HasSuffix(basePath, PathObjectBasePath) {
		basePath += PathObjectBasePath
	}

	// Initialize the internal secret value from an environment variable.
	// Note: This is important and secure because it resides deep within the binary internals and should not be left unset in production.
	internalSecretValue = os.Getenv(INTERNAL_SECRET_VALUE)
	if internalSecretValue == "" {
		panic(constant.InternelSecretEnvContextLog)
	}
}

// RegisterHandlersGin registers the HTTP handlers for the URL shortener service using the Gin
// web framework. It sets up the routes for retrieving, creating, and updating shortened URLs.
// The InternalOnly middleware is applied to the POST and PUT routes to protect them from public access.
func RegisterHandlersGin(router *gin.Engine, datastoreClient *datastore.Client) {
	// Register handlers with the custom or default base path.
	// For example, if CUSTOM_BASE_PATH is "/api/", the GET route will be "/api/:id",
	// the POST route will be "/api/", and the PUT route will be "/api/:id".
	router.GET(basePath+PathObjectID, getURLHandlerGin(datastoreClient))
	router.POST(basePath, InternalOnly(), postURLHandlerGin(datastoreClient))
	router.PUT(basePath+PathObjectID, InternalOnly(), editURLHandlerGin(datastoreClient))      // New PUT route for editing URLs
	router.DELETE(basePath+PathObjectID, InternalOnly(), deleteURLHandlerGin(datastoreClient)) // New DELETE route for deleting URLs
}

// generateShortID generates a unique short identifier for a URL.
//
// This function attempts to generate a unique short ID suitable for use in the datastore.
// It retries until a unique ID is found. If it cannot generate a unique ID after a predefined
// number of attempts, it returns an error, potentially indicating an issue with the underlying
// system or collision space.
//
// Note: This function is specifically tailored for generating unique short IDs for the datastore
// and may be adapted for other purposes in the future.
func generateShortID(ctx context.Context, dsClient *datastore.Client) (string, error) {
	id, err := shortid.GenerateUniqueDataStore(ctx, dsClient, 5)
	if err != nil {
		return "", err // If there's an error generating the ID, return it immediately.
	}
	return id, nil // If the ID is unique, return it.
}

// NewRateLimiter creates a new rate limiter for a client if it doesn't exist, or returns the existing one.
func NewRateLimiter(key string, r rate.Limit, b int) *rate.Limiter {
	RateLimiterStore.Lock()
	defer RateLimiterStore.Unlock()

	limiter, exists := RateLimiterStore.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(r, b)
		RateLimiterStore.limiters[key] = limiter
	}

	return limiter
}
