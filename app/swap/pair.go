package swap

import (
	"errors"
	"math/big"
)

var ErrInsufficientReserve = errors.New("insufficient reserve")

type PairTrade struct {
	Token0 TokenAmount
	Token1 TokenAmount
}

var (
	c998  = big.NewInt(998)
	c1000 = big.NewInt(1000)
)

func NewPair(tokenAmountA TokenAmount, tokenAmountB TokenAmount) *PairTrade {
	return &PairTrade{
		Token0: tokenAmountA,
		Token1: tokenAmountB,
	}
}

func (p PairTrade) GetOutputAmount(inputAmount TokenAmount) (TokenAmount, error) {
	if p.getReserve0() <= 0 || p.getReserve1() <= 0 {
		return TokenAmount{}, ErrInsufficientReserve
	}

	inputReserve := p.getReserveOf(inputAmount.Token)
	outputReserve := p.Token0
	if p.Token0.Token.IsEqual(inputAmount.Token) {
		outputReserve = p.Token1
	}

	inputAmountWithFee := inputAmount.Amount * 998
	numerator := inputAmountWithFee * outputReserve.Amount
	denominator := (inputReserve.Amount * 1000) + inputAmountWithFee

	outputAmount := TokenAmount{
		Token:  outputReserve.Token,
		Amount: numerator / denominator,
	}

	return outputAmount, nil
}

func (p PairTrade) GetOutputAmountInPip(inputAmount TokenAmount) (TokenAmount, error) {
	inputReserve := p.getReserveOf(inputAmount.Token)
	outputReserve := p.Token0
	if p.Token0.Token.IsEqual(inputAmount.Token) {
		outputReserve = p.Token1
	}

	inputAmountWithFee := new(big.Int).Mul(inputAmount.AmountInPip, c998)
	numerator := new(big.Int).Mul(inputAmountWithFee, outputReserve.AmountInPip)
	denominator := new(big.Int).Add(new(big.Int).Mul(inputReserve.AmountInPip, c1000), inputAmountWithFee)

	outputAmount := TokenAmount{
		Token:       outputReserve.Token,
		AmountInPip: numerator.Quo(numerator, denominator),
	}

	return outputAmount, nil
}

func (p PairTrade) GetInputAmount(outputAmount TokenAmount) (TokenAmount, error) {
	if p.getReserve0() == 0 || p.getReserve1() == 0 || p.getReserveOf(outputAmount.Token).Amount < outputAmount.Amount {
		return TokenAmount{}, ErrInsufficientReserve
	}

	outputReserve := p.getReserveOf(outputAmount.Token)
	inputReserve := p.Token0
	if p.Token0.Token.IsEqual(outputAmount.Token) {
		inputReserve = p.Token1
	}

	numerator := inputReserve.Amount * outputAmount.Amount * 1000
	denominator := (outputReserve.Amount - outputAmount.Amount) * 998

	amount := 0.0
	if denominator != 0 {
		amount = (numerator / denominator) + 0.000000000000000001
	}

	return TokenAmount{
		Token:  inputReserve.Token,
		Amount: amount,
	}, nil
}

func (p PairTrade) GetInputAmountInPip(outputAmount TokenAmount) (TokenAmount, error) {
	outputReserve := p.getReserveOf(outputAmount.Token)
	inputReserve := p.Token0
	if p.Token0.Token.IsEqual(outputAmount.Token) {
		inputReserve = p.Token1
	}

	numerator := new(big.Int).Mul(new(big.Int).Mul(inputReserve.AmountInPip, outputAmount.AmountInPip), c1000)
	denominator := new(big.Int).Mul(new(big.Int).Sub(outputReserve.AmountInPip, outputAmount.AmountInPip), c998)

	amount := big.NewInt(0)
	if denominator.Cmp(amount) != 0 {
		amount = new(big.Int).Add(new(big.Int).Div(numerator, denominator), big.NewInt(1))
	}

	return TokenAmount{
		Token:       inputReserve.Token,
		AmountInPip: amount,
	}, nil
}

func (p PairTrade) getReserve0() float64 {
	return p.Token0.Amount
}

func (p PairTrade) getReserve1() float64 {
	return p.Token1.Amount
}

func (p PairTrade) getReserveOf(token Token) TokenAmount {
	if p.Token0.Token.IsEqual(token) {
		return p.Token0
	}

	return p.Token1
}
