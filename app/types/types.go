package types

import (
	"github.com/MinterTeam/swap-router/app/swap"
)

type TradeSearch struct {
	FromCoinId uint64
	ToCoinId   uint64
	TradeType  swap.TradeType
	Amount     string
	Trade      chan *swap.Trade
}
