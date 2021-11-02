package services

import (
	"errors"
	"github.com/MinterTeam/swap-router/app/config"
	"github.com/MinterTeam/swap-router/app/swap"
	"github.com/MinterTeam/swap-router/app/types"
	"math/big"
)

type Swap struct {
	poolService *Pool

	tradeSearchJobs chan types.TradeSearch
}

func NewSwapService(cfg config.WorkersConfig, ps *Pool) *Swap {
	service := &Swap{
		poolService:     ps,
		tradeSearchJobs: make(chan types.TradeSearch),
	}

	go service.runWorkers(cfg.FindRouteWorkersCount)
	return service
}

func (s *Swap) runWorkers(workersCount int) {
	for w := 1; w <= workersCount; w++ {
		go s.runFindRouteWorker(s.tradeSearchJobs)
	}
}

func (s *Swap) runFindRouteWorker(jobs <-chan types.TradeSearch) {
	//for j := range jobs {
	//	trade, _ := s.findRoute(j.FromCoinId, j.ToCoinId, j.TradeType, j.Amount)
	//	j.Trade <- trade
	//}
}

func (s *Swap) findRoute(fromCoinId uint64, toCoinId uint64, tradeType swap.TradeType, amount *big.Int) (*swap.Trade, error) {
	ts := types.TradeSearch{
		FromCoinId: fromCoinId,
		ToCoinId:   toCoinId,
		TradeType:  tradeType,
		Amount:     amount,
		Trade:      make(chan *swap.Trade),
	}

	s.tradeSearchJobs <- ts
	trade := <-ts.Trade

	if trade == nil {
		return nil, errors.New("path not found")
	}

	return trade, nil
}

func (s *Swap) FindRoute(fromCoinId uint64, toCoinId uint64, tradeType swap.TradeType, amount string) (trade *swap.Trade, err error) {
	pairs, trade := s.poolService.GetTradePairs(), &swap.Trade{}
	if tradeType == swap.TradeTypeExactInput {
		trade, err = swap.GetBestTradeExactIn(pairs, swap.NewToken(toCoinId), swap.NewTokenAmount(swap.NewToken(fromCoinId), amount), 4)
	} else {
		trade, err = swap.GetBestTradeExactOut(pairs, swap.NewToken(fromCoinId), swap.NewTokenAmount(swap.NewToken(toCoinId), amount), 4)
	}

	if err != nil {
		return nil, err
	}

	if trade == nil {
		return nil, errors.New("path not found")
	}

	return trade, nil
}
