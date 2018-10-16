package hdkey

import (
	"fmt"
	"testing"
	"errors"
	"reflect"
	"encoding/hex"
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
	seed1, err1 := hex.DecodeString("000102030405060708090a0b0c0d0e0f") 
	seed2, err2 := hex.DecodeString("fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542") 

	if err1 != nil || err2 != nil {
		t.Errorf("Fail to decode string to byte array. detail:%v, %v", err1, err2)
	}
	tests := []struct {
		name	string
		seed	[]byte
		extpub string
		extprv	string
		errors	error
	}{
		{ name: "Test Vector1", seed: seed1, extprv: "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi" },
		{ name: "Test Vector2", seed: seed2, extprv: "xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U" },
	}
	for _, test := range tests {
		fmt.Printf("TEST: %s, seed: %x\n", test.name, test.seed)
		masterkey, e1 := MasterGen(test.seed)
		serialized, e2 := masterkey.Serialize()
		if e1 != nil {
			t.Errorf("Fail to generate master key. detail:%v\n", e1)
		}
		if e2 != nil {
			t.Errorf("Fail to serialize generated key. detail:%v\n", e2)
		}
		if serialized != test.extprv {
			t.Errorf("Wrong key. your key:%v\n", serialized)
		}
	}
}
