package shortid

import (
	"crypto/rand"
	"encoding/base64"
)

// Generate creates a URL-friendly short ID.
func Generate(length int) (string, error) {
	// Adjust length to account for base64 encoding
	encodedLength := length * 3 / 4

	// Generate random bytes
	randomBytes := make([]byte, encodedLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// URL-friendly encoding
	encoded := base64.URLEncoding.EncodeToString(randomBytes)

	// Trim padding
	encoded = encoded[:length]

	return encoded, nil
}
