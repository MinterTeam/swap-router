package resources

import (
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type Coin struct {
	ID     uint32 `json:"id"`
	Symbol string `json:"symbol"`
}

func (Coin) Transform(model ItemInterface, params ...ParamInterface) Interface {
	coin := model.(models.Coin)
	return Coin{
		ID:     uint32(coin.ID),
		Symbol: coin.GetSymbol(),
	}
}
