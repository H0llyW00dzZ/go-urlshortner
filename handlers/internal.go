package handlers

import (
	"net/http"
	"net/url"

	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"github.com/gin-gonic/gin"
)

// isValidURL checks if the URL is in a valid format.
func isValidURL(urlStr string) bool {
	u, err := url.ParseRequestURI(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// InternalOnly creates a middleware that restricts access to a route to internal services only.
// It checks for a specific header containing a secret value that should match an environment
// variable to allow the request to proceed. If the secret does not match or is not provided,
// the request is aborted with a 403 Forbidden status.
//
// Additionally, this middleware enforces rate limiting to prevent abuse.
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

		// Check if the request is allowed by the rate limiter.
		if !applyRateLimit(c) {
			// If the rate limit is exceeded, we should not continue processing the request.
			return
		}

		// If the header matches and the rate limiter allows it, proceed with the request.
		c.Next()
	}
}
