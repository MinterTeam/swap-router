package api

import (
	"github.com/MinterTeam/swap-router/app/helpers"
	"github.com/MinterTeam/swap-router/app/swap"
	"math/big"
)

type FindSwapPoolRouteRequest struct {
	Coin0 string `uri:"coin0"  binding:"required"`
	Coin1 string `uri:"coin1"  binding:"required"`
}

type FindSwapPoolRouteRequestQuery struct {
	Amount    string `form:"amount" binding:"required,numeric"`
	TradeType string `form:"type"   binding:"required,oneof=input output"`
}

func (req *FindSwapPoolRouteRequestQuery) GetAmount() *big.Int {
	return helpers.Str2BigInt(req.Amount)
}

func (req *FindSwapPoolRouteRequestQuery) GetTradeType() swap.TradeType {
	return swap.Str2TradeType(req.TradeType)
}
