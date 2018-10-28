package utils

import (
	"math/big"
)

var (
	bitcoinBase58Strings = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

func reverse(b []byte) []byte {
	// reverse the order
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func Encode(b []byte) string {
	num := new(big.Int).SetBytes(b)
	radix := big.NewInt(58)

	ret := make([]byte, 0)
	for num.Sign() == 1 {
		mod := new(big.Int)
		num.DivMod(num, radix, mod)
		ret = append(ret, bitcoinBase58Strings[mod.Int64()])
	}
	ret = reverse(ret)
	return string(ret)
}

func Decode(s string) []byte {
	revstr := string(reverse([]byte(s)))
	ret := big.NewInt(0)
	radix := big.NewInt(1)
	fiftyEight := big.NewInt(58)
	for _, b := range revstr {
		radix2 := new(big.Int).Set(radix)
		for j, c := range bitcoinBase58Strings {
			if b == c {
				ret.Add(ret, radix2.Mul(radix2, big.NewInt(int64(j))))
			}
		}
		radix.Mul(radix, fiftyEight)
	}
	return ret.Bytes()
}
