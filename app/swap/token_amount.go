package swap

import "math/big"

type TokenAmount struct {
	Token  Token
	Amount *big.Int
}

func NewTokenAmount(token Token, amount *big.Int) TokenAmount {
	return TokenAmount{Token: token, Amount: amount}
}

func (ta TokenAmount) GetAmount() *big.Int {
	return ta.Amount
}

func (ta TokenAmount) GetCurrency() Token {
	return ta.Token
}
