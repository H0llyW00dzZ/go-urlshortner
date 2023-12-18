// Gopher Unit Testing was here
package shortid

import (
	"testing"
)

// TestGenerate ensures that the Generate function returns a string of the correct length and no error.
func TestGenerate(t *testing.T) {
	length := 1337
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
	length := 1337
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

// TestGenerate_ExpectedLength tests that the Generate function returns a string of the expected length.
func TestGenerate_ExpectedLength(t *testing.T) {
	lengths := []int{6, 8, 10} // More realistic lengths for a short ID
	for _, length := range lengths {
		id, err := Generate(length)
		if err != nil {
			t.Fatalf("Generate(%d) returned an unexpected error: %v", length, err)
		}
		if got := len(id); got != length {
			t.Errorf("Generate(%d) returned a string of length %d, want %d", length, got, length)
		}
	}
}

// TestGenerate_Uniqueness checks if the Generate function returns unique values on subsequent calls.
func TestGenerate_Uniqueness(t *testing.T) {
	length := 5       // A more realistic length for a short ID
	iterations := 100 // Reduce the number of iterations
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

// TestGenerate_BoundaryConditions tests that the Generate function returns an error for invalid lengths.
func TestGenerate_BoundaryConditions(t *testing.T) {
	testCases := []int{0, -1}
	for _, length := range testCases {
		_, err := Generate(length)
		if err == nil {
			t.Errorf("Generate(%d) should return an error", length)
		}
	}
}
