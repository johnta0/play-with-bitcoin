package hdkey

// This is  an implementation of https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki

import (
	"encoding/hex"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"

	// TODO: Impl it myself in the future
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcd/btcec"
)

const (
	// MinSeedBytes defines min value of seed length in bytes
	MinSeedBytes = 16 // 128 bits
	// MaxSeedBytes defines max value of seed length in bytes
	MaxSeedBytes = 64 // 512 bits
	// RecommendedBytes is the recommended seed length in bytes
	RecommendedBytes = 32 // 256 bits
)

var (
	// ErrInvalidSeedLength describes an error in which provided seed length
	// is not in the specified range
	ErrInvalidSeedLength = fmt.Errorf("Seed length must be between %d and %d bits", MinSeedBytes*8, MaxSeedBytes*8)
	ErrInvalidSeedValue  = fmt.Errorf("Invalid Seed. Please try another seed")
	// versions
	MainPub, _ = hex.DecodeString("0488B21E")
	MainPrv, _ = hex.DecodeString("0488ADE4")
	TestPub, _ = hex.DecodeString("043587CF")
	TestPrv, _ = hex.DecodeString("04358394")

	// ErrDeriveHardenedFromPub describes an error which is throwed when
	// you try to derive hardened child key from extended parent public key
	ErrDeriveHardenedFromPub = fmt.Errorf("Cannot derive hardened child key from extended public key.")

	// version bytes
	MainPubVer, _ = hex.DecodeString("0488B21E")
	MainPrvVer, _ = hex.DecodeString("0488ADE4")
	TestPubVer, _ = hex.DecodeString("043587CF")
	TestPrvVer, _ = hex.DecodeString("04358394")
)

// ExtKey type houses params for extended private key
type ExtKey struct {
	key               []byte // 33 bytes
	pubkey            []byte
	chainCode         []byte // 32 bytes
	version           []byte // 4 byte
	depth             uint8  // 1 byte
	parentFingerPrint []byte
	childNum          uint32 // 4bytes
	isPrivate         bool   // true => privkey, false => pubkey
}

// NewExtKey returns a new instnace of ExtKey
func NewExtKey(key []byte, chainCode []byte, version []byte, depth uint8,
	parentFingerPrint []byte, childNum uint32, isPrivate bool) *ExtKey {
	return &ExtKey{
		key:               key,
		chainCode:         chainCode,
		version:           version,
		depth:             depth,
		parentFingerPrint: parentFingerPrint,
		childNum:          childNum,
		isPrivate:         isPrivate,
	}
}

// MasterGen returns master key derived from seed.
//
// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#Master_key_generation
func MasterGen(seed []byte) (*ExtKey, error) {
	// I = HMAC-SHA512(Key = "Bitcoin seed", Data = S)
	mac := hmac.New(sha512.New, []byte("Bitcoin seed"))
	mac.Write(seed)
	iAll := mac.Sum(nil)

	iR := iAll[32:] // chainCode
	iL := iAll[:32] // privkey
	privkey := new(big.Int).SetBytes(iL)
	n := btcec.S256().N
	if privkey.Sign() == 0 || privkey.Cmp(n) == 1 {
		return nil, ErrInvalidSeedValue
	}
	return NewExtKey(
		iL, // key
		iR, // chainCode
		// TODO: Be able to make choice main/testnet
		MainPrv,
		0,               // depth
		make([]byte, 4), // parentFP is 0x00000000 if masterkey
		0,               //childNum
		true,
	), nil
}

// this function serialize ExtKey.
//
// Serialization Format:
// version || depth || parentFP || childNum || chainCode || 0x00 || key || checksum
func (k *ExtKey) Serialize() (string, error) {
	ret := make([]byte, 0)
	ret = append(ret, k.version...)
	ret = append(ret, k.depth)
	ret = append(ret, k.parentFingerPrint...)
	childNum := make([]byte, 4)
	binary.BigEndian.PutUint32(childNum, k.childNum)
	ret = append(ret, childNum...)
	ret = append(ret, k.chainCode...)
	if k.isPrivate {
		ret = append(ret, 0x00)
		ret = append(ret, k.key...)
	} else {
		// unimplemented
	}
	// checksum = sha256(sha256(ret))
	hash := sha256.Sum256(ret)
	doublehash := sha256.Sum256(hash[:])
	checksum := doublehash[:][:4]

	ret = append(ret, checksum...)

	return base58.Encode(ret), nil
}

// SeedGen returns seed.
//
// Generate a seed byte sequence S of a chosen length (between 128 and 512 bits; 256 bits is advised)
// [16, 64] bytes, 32 bits advised
func SeedGen(length uint8) ([]byte, error) {
	// The seed range confining
	if length < MinSeedBytes || length > MaxSeedBytes {
		return nil, ErrInvalidSeedLength
	}

	seed := make([]byte, length)
	_, err := rand.Read(seed) // throw length away
	if err != nil {
		return nil, err
	}
	return seed, nil
}

// DerivePubkey returns pubkey bytes derived from private key,
// if ExtKey instnace isPrivate
func (k *ExtKey) DerivePubkey() []byte {
	if !k.isPrivate {
		return k.key
	}

	if k.pubkey == nil {
		x, y := btcec.S256().ScalarBaseMult(k.key)
		pk := btcec.PublicKey{btcec.S256(), x, y}
		k.pubkey = pk.SerializeCompressed()
	}
	return k.pubkey
}

func (k *ExtKey) getPubkeyBytes() []byte {
	if !k.isPrivate {
		return k.pubkey
	}

	if k.pubkey == nil {
		_, pub := btcec.PrivKeyFromBytes(btcec.S256(), k.key)
		k.pubkey = pub.SerializeCompressed()
	}
	return k.pubkey
}

// DeriveChildKey returns derived childed extended key
func (k *ExtKey) DeriveChildKey(i uint32) (*ExtKey, error) {

	// There are four cases:
	// 	1. Private extkey => Private child hardend key
	//	2. Private extkey => Private child normal key
	//	3. Public extkey => Public child hardened key (return invalid)
	//	4. Public extkey => Public child normal key

	// the child extkey is hardened or not
	toBeHardened := i >= 0x80000000 // 2^31 = 8 * 16^7

	// Case 3 is invalid. Return error early
	if !k.isPrivate && toBeHardened {
		return nil, ErrDeriveHardenedFromPub
	}

	// When hardened: Data = 0x00 || ser256(kpar) || ser32(i)
	// When not:	  Data = serP(point(kpar))
	data := make([]byte, 37)
	binary.BigEndian.PutUint32(data[33:], i)
	if toBeHardened {
		// case 1
		copy(data[1:], k.key)
	} else {
		// case 2, 4
		// I = HMAC-SHA512(Key = cpar, Data = serP(point(kpar) or K_par) || ser32(i)).
		copy(data, k.getPubkeyBytes())
	}

	mac512 := hmac.New(sha512.New, k.chainCode) // use chain code as the key
	mac512.Write(data)
	// I = HMAC-SHA512(c_par, data), I_L => be used to derive child, I_R => chain code
	iAll := mac512.Sum(nil)

	iL := iAll[:len(iAll)/2]
	childChainCode := iAll[len(iAll)/2:]

	return nil, NewExtKey()
}

// DerivePubkey returns public key derived from given private key
// func (k *ExtKey) DerivePubkey(privkey []byte) []byte {
// 	return
// }

// Neuter returns extended public key corresponding to given
// extended private key
func (k *ExtKey) Neuter() *ExtKey {
	if !k.isPrivate {
		return k
	}

	// version: from priv to pub
	// k.key: from priv to pub
	// k.pubkey: from pub to nil

	var version []byte
	if reflect.DeepEqual(k.version, MainPrvVer) {
		version = MainPubVer
	} else if reflect.DeepEqual(k.version, TestPrvVer) {
		version = TestPubVer
	}
	

	return NewExtKey(
		k.getPubkeyBytes(),
		k.chainCode,
		version,
		k.depth,
		k.parentFingerPrint,
		k.childNum,
		false,
	)
}
