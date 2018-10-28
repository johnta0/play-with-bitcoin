package utils

import (
	"testing"
	"reflect"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		value  string
		err     error
	}{
		// valid ones
		{name: "case1: 123456789", input: []byte("123456789"), value: "dKYWwnRHc7Ck"},
		{name: "case2: bitcoin", input: []byte("bitcoin"), value: "4jJc4sAwPs"},
	}
	for _, test := range tests {
		encoded := Encode(test.input)
		if encoded != test.value {
			t.Errorf("value mismatch. output: %v, correct result: %v", encoded, test.value)
		}
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		value  []byte
		err     error
	}{
		// valid ones
		{name: "case1: 123456789", input: "dKYWwnRHc7Ck", value: []byte("123456789")},
		{name: "case2: bitcoin", input: "4jJc4sAwPs", value: []byte("bitcoin")},
	}
	for _, test := range tests {
		decoded := Decode(test.input)
		if !reflect.DeepEqual(decoded, test.value) {
			t.Errorf("value mismatch. output: %v, correct value: %v", decoded, test.value)
		}
	}
}
