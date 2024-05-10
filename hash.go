package cuid2

import (
	"crypto/sha512"
	"math/big"
)

func bufToInt64(buf []byte) *big.Int {
	var bits uint = 8
	value := big.NewInt(0)
	for _, v := range buf {
		value = value.Lsh(value, bits).Add(value, big.NewInt(int64(v)))
	}
	return value
}

func hash(input string) string {
	hasher := sha512.New()
	hasher.Write([]byte(input))
	bytes := hasher.Sum(nil)
	num := bufToInt64(bytes)
	return base36EncodeBigint(num)[1:]
}
