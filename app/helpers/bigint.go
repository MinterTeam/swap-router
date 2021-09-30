package helpers

import "math/big"

// default amount of pips in 1 bip
var pipInBip = big.NewFloat(1000000000000000000)

func Str2BigInt(value string) *big.Int {
	v, _ := new(big.Int).SetString(value, 10)
	return v
}

func PipStr2Bip(value string) float64 {
	if value == "" {
		value = "0"
	}

	floatValue, _ := new(big.Float).SetPrec(500).SetString(value)
	f, _ := new(big.Float).SetPrec(500).Quo(floatValue, big.NewFloat(1000000000000000000)).Float64()
	return f
}

func BipFloatToStr(value float64) string {
	return big.NewFloat(value).SetMode(big.ToZero).Text('f', 12)
}
