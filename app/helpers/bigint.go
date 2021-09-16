package helpers

import "math/big"

func Str2BigInt(value string) *big.Int {
	v, _ := new(big.Int).SetString(value, 10)
	return v
}
