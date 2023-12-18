package handlers

import (
	"fmt"
	"strings"

	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"github.com/gin-gonic/gin"
)

// constructFullShortenedURL constructs the full shortened URL from the request and the base path.
func constructFullShortenedURL(c *gin.Context, id string) string {
	// Check for the X-Forwarded-Proto header to determine the scheme.
	scheme := c.GetHeader(constant.HeaderXProto)
	if scheme == "" {
		// Fallback to checking the TLS property of the request if the header is not set.
		if c.Request.TLS != nil {
			scheme = constant.HeaderSchemeHTTPS
		} else {
			scheme = constant.HeaderSchemeHTTP
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
