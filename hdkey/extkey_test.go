package hdkey

import (
	"testing"
)

// TestSeedGen ensures that SeedGen returns error when specified length is not
// between 128 and 512 bits.
func TestSeedGen(t *testing.T) {
	tests := []struct {
		length	uint
		isValid	bool
		err		error
	} {
		// valid ones
		{ length: "128", isValid: true },
		// invalid ones
		{ length: "7", isValid: false },
	}
	for _, test := range tests {
		seed, err := SeedGen(test.length)
		if test.isValid && err != nil {
			t.Error("The length %d should be invalid, but the function didn't return an error.",
					+ test.length)
		}
		if !test.isValid && err == nil {
			t.Error("Test length %d should be valid.", test.length)
		}
	}
}
