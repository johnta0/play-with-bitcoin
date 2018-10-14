package hdkey

// This is  an implementation of https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki

import (
	"fmt"
	"crypto/rand"
)

const(
	// Seed: between 128 and 512 bits; 256 bits is advised
	MinBytes = 16 // 128 bits
	MaxBytes = 64 // 512 bits
	RecommendedBytes = 32 // 256 bits
)

var(
	ErrInvalidSeedLength = fmt.Errorf("Seed length must be between %d and %d bits"
		, MinSeedBytes*8, MaxSeedBytes*8)
)

// ExtendedKey type houses params for extended private key
type ExtKey struct {
	key []byte // 33 bytes
	chainCode []byte // 32 bytes
	version []byte // 4 byte
	depth uint8 // 1 byte
	parentFingerPrint []byte
	childNum uint32 // 4bytes
	isPrivate bool // true => privkey, false => pubkey
}

// MasrterGen return master key derived from seed.
//
// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#Master_key_generation
func (k *ExtKey)MasterGen(seed []byte) (*ExtKey, error) {
	if len(seed) < MinSeedBytes || len(seed) > MaxSeedBytes {
		return nil, ErrInvalidSeedLength
	}

	return NewExtKey(key, 0, true)
}

// SeedGen returns seed.
func SeedGen(length uint) ([]byte, error) {
	if length < MinBytes || lenght > MaxBytes {
		return nil, ErrInvalidSeedLength
	}

	seed := make([]byte, length)
	_, err := rand.Read(seed) // throw length away
	if err != nil {
		return nil, err
	}
	return seed
}

// ChildKeyDeriv returns derived childed key by index
func (k *ExtKey) DeriveChildKey(index uint) (*ExtPrivKey) {
	return
}

// DerivePubkey returns public key derived from given private key
func (k *ExtKey) DerivePubkey(privkey []byte) []byte {
	return
}
