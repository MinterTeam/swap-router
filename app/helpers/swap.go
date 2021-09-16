package helpers

import "github.com/MinterTeam/swap-router/app/swap"

func Str2TradeType(tradeType string) swap.TradeType {
	if tradeType == "input" {
		return swap.TradeTypeExactInput
	}

	return swap.TradeTypeExactOutput
}
