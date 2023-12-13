package shortid

import (
	"crypto/rand"
	"encoding/base64"
)

// Generate creates a URL-friendly short ID.
func Generate(length int) (string, error) {
	// Generate a buffer larger than needed to ensure we have enough characters
	// after base64 encoding to satisfy the requested length.
	bufferSize := length * 3 / 4
	if length%3 != 0 {
		bufferSize++ // Compensate for partial encoding groups
	}

	// Generate random bytes
	randomBytes := make([]byte, bufferSize)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// URL-friendly encoding
	encoded := base64.URLEncoding.EncodeToString(randomBytes)

	// Ensure we have enough encoded characters to satisfy the requested length
	for len(encoded) < length {
		extraBytes := make([]byte, 1)
		if _, err := rand.Read(extraBytes); err != nil {
			return "", err
		}
		encoded += base64.URLEncoding.EncodeToString(extraBytes)
	}

	// Trim to the requested length and remove any base64 padding
	encoded = encoded[:length]
	encoded = string([]rune(encoded)) // Convert to slice of runes to handle multi-byte characters
	return encoded, nil
}
