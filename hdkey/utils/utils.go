package utils

import (
	"math/big"
)

var (
	BitcoinBase58Strings = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

// reverse the order
func reverse(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

// Encode returns string encoded by base58
// arg: byte array
func Encode(b []byte) string {
	num := new(big.Int).SetBytes(b)
	radix := big.NewInt(58)

	ret := make([]byte, 0)
	// Divide num by 58 and append char corresponded to its remainder until num is less than 58
	for num.Sign() == 1 {
		mod := new(big.Int)
		num.DivMod(num, radix, mod)
		ret = append(ret, BitcoinBase58Strings[mod.Int64()])
	}
	ret = reverse(ret)
	return string(ret)
}

// Decode returns byte array decoded by base58
// arg: string
func Decode(s string) []byte {
	revstr := string(reverse([]byte(s)))
	ret := big.NewInt(0)
	radix := big.NewInt(1)
	fiftyEight := big.NewInt(58)
	// Convert base58string to bigint
	for _, b := range revstr {
		radix2 := new(big.Int).Set(radix)
		for j, c := range BitcoinBase58Strings {
			if b == c {
				ret.Add(ret, radix2.Mul(radix2, big.NewInt(int64(j))))
			}
		}
		radix.Mul(radix, fiftyEight)
	}
	return ret.Bytes()
}
