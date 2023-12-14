// Package shortid provides functionality to generate cryptographically secure,
// URL-friendly unique identifiers.
//
// Copyright (c) 2023 H0llyW00dzZ
package shortid

import (
	"crypto/rand"
	"encoding/base64"
)

// Generate creates a cryptographically secure, URL-friendly short ID of a specified length.
// The ID is generated using random bytes, which are then base64 URL encoded to ensure
// they can be safely used in URLs. The length parameter specifies the desired length
// of the final encoded ID. If the length is not a multiple of 4, the function compensates
// to ensure the final ID has the correct length.
//
// The function returns the generated short ID or an error if the random byte generation fails.
//
// Example usage:
//
//	id, err := shortid.Generate(10)
//	if err != nil {
//	    log.Fatalf("Failed to generate short ID: %v", err)
//	}
//	fmt.Println("Generated short ID:", id)
func Generate(length int) (string, error) {
	// Calculate the buffer size needed to ensure the base64 encoded string meets the requested length.
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
