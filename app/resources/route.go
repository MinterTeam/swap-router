package resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
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
		AmountIn:  helpers.Pip2BipStr(trade.InputAmount.GetAmountInPip()),
		AmountOut: helpers.Pip2BipStr(trade.OutputAmount.GetAmountInPip()),
		Coins:     TransformCollection(route, new(Coin)),
	}
}
