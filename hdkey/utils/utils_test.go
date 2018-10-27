package utils

import (
	"testing"
)

func TestEncode(t *testing.T) {
	encoded := Encode([]byte("123456789"))
	if encoded != "dKYWwnRHc7Ck" {
		t.Errorf("value mismatch. %v", encoded)
	}
}
