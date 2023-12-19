package shortid

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/H0llyW00dzZ/go-urlshortner/datastore"
)

// Generate creates a cryptographically secure, URL-friendly short ID of a specified length.
func Generate(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be a positive integer, got %d", length)
	}
	return generateRandomString(length)
}

// GenerateUniqueDataStore creates a unique, cryptographically secure, URL-friendly short ID.
//
// Note: This function has been renamed to GenerateUniqueDataStore to avoid confusion with similarly named 'Generate' functions for other databases in the future.
func GenerateUniqueDataStore(ctx context.Context, client *datastore.Client, length int) (string, error) {
	const maxRetries = 1337 // Maximum number of retries to find a unique ID
	for i := 0; i < maxRetries; i++ {
		id, err := generateRandomString(length)
		if err != nil {
			return "", fmt.Errorf("error generating random string: %w", err)
		}

		// Check if the ID already exists in the datastore.
		_, err = datastore.GetURL(ctx, client, id)
		if err != nil {
			if errors.Is(err, datastore.ErrNotFound) {
				// The ID does not exist, so it is unique
				return id, nil
			}
			// Some other error occurred when checking the ID
			return "", fmt.Errorf("error checking ID uniqueness: %w", err)
		}
		// The ID exists, so it is not unique, try generating another one
	}

	return "", errors.New("failed to generate a unique short ID after several attempts")
}

// generateRandomString generates a random, URL-friendly string of a specified length.
func generateRandomString(length int) (string, error) {
	bufferSize := length * 3 / 4
	if length%3 != 0 {
		bufferSize++ // Compensate for partial encoding groups
	}

	// Generate a slice of random bytes.
	randomBytes := make([]byte, bufferSize)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// Encode the random bytes into a URL-friendly base64 string.
	encoded := base64.URLEncoding.EncodeToString(randomBytes)

	// If the encoded string is shorter than the requested length, append more random characters.
	// This loop is inefficient and can be improved.
	for len(encoded) < length {
		extraBytesNeeded := length - len(encoded)
		extraBytes := make([]byte, (extraBytesNeeded*3+3)/4) // Adjusted buffer size calculation
		if _, err := rand.Read(extraBytes); err != nil {
			return "", err
		}
		encoded += base64.URLEncoding.EncodeToString(extraBytes)
	}

	// Trim the encoded string to the requested length and remove any base64 padding.
	encoded = encoded[:length]
	encoded = string([]rune(encoded)) // Convert to a slice of runes to handle multi-byte characters.
	return encoded, nil
}
