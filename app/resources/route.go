package resources

import (
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/swap-router/app/helpers"
	"github.com/MinterTeam/swap-router/app/swap"
)

type Route struct {
	SwapType  string      `json:"swap_type"`
	AmountIn  string      `json:"amount_in"`
	AmountOut string      `json:"amount_out"`
	Coins     []Interface `json:"coins"`
}

func (r Route) Transform(route []models.Coin, trade *swap.Trade) Route {
	return Route{
		SwapType:  "pool",
		AmountIn:  helpers.BipFloatToStr(trade.InputAmount.GetAmount()),
		AmountOut: helpers.BipFloatToStr(trade.OutputAmount.GetAmount()),
		Coins:     TransformCollection(route, new(Coin)),
	}
}
