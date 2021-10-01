package swap

import (
	"github.com/MinterTeam/swap-router/app/helpers"
	"math/big"
)

type TokenAmount struct {
	Token       Token
	Amount      float64
	AmountInPip *big.Int
}

func NewTokenAmount(token Token, amount string) TokenAmount {
	return TokenAmount{
		Token:       token,
		Amount:      helpers.PipStr2Bip(amount),
		AmountInPip: helpers.Str2BigInt(amount),
	}
}

func (ta TokenAmount) GetAmount() float64 {
	return ta.Amount
}

func (ta TokenAmount) GetAmountInPip() *big.Int {
	return ta.AmountInPip
}

func (ta TokenAmount) GetCurrency() Token {
	return ta.Token
}
