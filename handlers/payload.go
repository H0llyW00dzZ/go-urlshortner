package handlers

import (
	"fmt"

	"github.com/H0llyW00dzZ/go-urlshortner/logmonitor/constant"
	"github.com/gin-gonic/gin"
)

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

// bindUpdatePayload binds the JSON payload to the UpdateURLPayload struct and validates the new URL format.
func bindUpdatePayload(c *gin.Context) (UpdateURLPayload, error) {
	var req UpdateURLPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, err
	}

	if req.NewURL == "" || !isValidURL(req.NewURL) {
		return req, fmt.Errorf(constant.InvalidNewURLFormatContextLog)
	}

	return req, nil
}
