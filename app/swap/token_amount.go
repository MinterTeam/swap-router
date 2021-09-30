package swap

type TokenAmount struct {
	Token  Token
	Amount float64
}

func NewTokenAmount(token Token, amount float64) TokenAmount {
	return TokenAmount{Token: token, Amount: amount}
}

func (ta TokenAmount) GetAmount() float64 {
	return ta.Amount
}

func (ta TokenAmount) GetCurrency() Token {
	return ta.Token
}
