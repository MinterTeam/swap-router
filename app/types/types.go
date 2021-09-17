package types

import (
	"github.com/MinterTeam/swap-router/app/swap"
	"math/big"
)

type TradeSearch struct {
	FromCoinId uint64
	ToCoinId   uint64
	TradeType  swap.TradeType
	Amount     *big.Int
	Trade      chan *swap.Trade
}
