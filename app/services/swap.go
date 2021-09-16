package services

import (
	"errors"
	"github.com/MinterTeam/swap-router/app/swap"
	"math/big"
)

type Swap struct {
	poolService *Pool
}

func NewSwapService(ps *Pool) *Swap {
	return &Swap{ps}
}

func (s *Swap) FindRoute(fromCoinId uint64, toCoinId uint64, tradeType swap.TradeType, amount *big.Int) (trade *swap.Trade, err error) {
	pairs, trades := s.poolService.GetTradePairs(), []swap.Trade{}
	if tradeType == swap.TradeTypeExactInput {
		trades, err = swap.GetBestTradeExactIn(pairs, swap.NewToken(toCoinId), swap.NewTokenAmount(swap.NewToken(fromCoinId), amount),
			swap.TradeOptions{MaxNumResults: 1, MaxHops: 4})
	} else {
		trades, err = swap.GetBestTradeExactOut(pairs, swap.NewToken(fromCoinId), swap.NewTokenAmount(swap.NewToken(toCoinId), amount),
			swap.TradeOptions{MaxNumResults: 1, MaxHops: 4})
	}

	if err != nil {
		return nil, err
	}

	if len(trades) == 0 {
		return nil, errors.New("path not found")
	}

	return &trades[0], nil
}
