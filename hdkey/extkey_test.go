package hdkey

import (
	"fmt"
	"testing"
<<<<<<< HEAD
	"errors"
	"reflect"
=======
	"encoding/hex"
>>>>>>> hdkey: impl mastergen but fail to test
)

// TestSeedGen ensures that SeedGen returns error when specified length is not
// between 128 and 512 bits.
func TestSeedGen(t *testing.T) {
	wantErr := errors.New("Seed length must be between 128 and 512 bits")

	tests := []struct {
		name	string
		length	uint8 // keep in mind that the unit is byte, not bit.
		isValid	bool
		err		error
	} {
		// valid ones
		{ name: "128 bits", length: 16, isValid: true},
		{ name: "256 bits (recommended)", length: 32, isValid: true },
		// invalid ones
		{ name: "56 bits (too short)", length: 7, isValid: false, err: wantErr },
	}

	for _, test := range tests {
		seed, err := SeedGen(test.length)
		if err == nil && len(seed) != int(test.length) {
			t.Errorf("test %s: length doesn't match: want %v, got %v", test.name, 
				test.length, len(seed))
			continue
		} else if !reflect.DeepEqual(err, test.err) {
			t.Errorf("test %s: unexpected error ------" +
				"want %v, got %v", test.name, test.err, err)
			continue
		}
	}
}

func TestMasterGen(t *testing.T) {
	seed, e := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	if e != nil {
		t.Error("Fail to decode string to byte array")
	}
	masterprivkey, err := MasterGen(seed)
	if err != nil {
		t.Error("Fail to generate masterkey")
	}
	fmt.Printf("+%v", masterprivkey)
}
