package swap

func Str2TradeType(tradeType string) TradeType {
	if tradeType == "input" {
		return TradeTypeExactInput
	}

	return TradeTypeExactOutput
}
