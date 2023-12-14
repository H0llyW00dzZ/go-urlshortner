// Gopher Unit Testing was here
package shortid

import (
	"testing"
)

// TestGenerate ensures that the Generate function returns a string of the correct length and no error.
func TestGenerate(t *testing.T) {
	length := 5
	id, err := Generate(length)
	if err != nil {
		t.Fatalf("Generate returned an unexpected error: %v", err)
	}
	if len(id) != length {
		t.Errorf("Generate returned a string of length %d, want %d", len(id), length)
	}
}

// TestGenerateUniqueness checks if the Generate function returns unique values on subsequent calls.
func TestGenerateUniqueness(t *testing.T) {
	length := 5
	iterations := 1000
	ids := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		id, err := Generate(length)
		if err != nil {
			t.Fatalf("Generate returned an unexpected error: %v", err)
		}
		if ids[id] {
			t.Fatalf("Generate returned a non-unique id: %v", id)
		}
		ids[id] = true
	}
}

// TestGenerateLengthVariance tests that the Generate function can handle different lengths.
func TestGenerateLengthVariance(t *testing.T) {
	for length := 1; length <= 10; length++ {
		id, err := Generate(length)
		if err != nil {
			t.Fatalf("Generate(%d) returned an unexpected error: %v", length, err)
		}
		if len(id) != length {
			t.Errorf("Generate(%d) returned a string of length %d, want %d", length, len(id), length)
		}
	}
}
