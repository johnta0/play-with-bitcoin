package utils

import (
	"math/big"
)

var (
	StringSet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
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
		ret = append(ret, StringSet[mod.Int64()])
	}
	ret = reverse(ret)
	return string(ret)
}

