package services

import (
	"errors"
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/swap-router/app/repositories"
	"sync"
)

type Coin struct {
	repository *repositories.Coin

	coins         sync.Map
	coinsBySymbol sync.Map
}

var (
	ErrNotFound = errors.New("not found")
)

func NewCoinService(r *repositories.Coin) *Coin {
	service := &Coin{repository: r}
	service.fillCoinsMap()
	return service
}

func (c *Coin) GetCoinById(id uint64) models.Coin {
	coin, _ := c.coins.Load(id)
	return coin.(models.Coin)
}

func (c *Coin) GetCoinIdBySymbol(symbol string) (uint64, error) {
	coin, ok := c.coinsBySymbol.Load(symbol)
	if ok {
		return uint64(coin.(models.Coin).ID), nil
	}

	return 0, ErrNotFound
}

func (c *Coin) ListenNewBlock(b blocks.Resource) {
	c.fillCoinsMap()
}

func (c *Coin) fillCoinsMap() {
	coins, _ := c.repository.GetAll()
	for _, coin := range coins {
		c.coins.Store(uint64(coin.ID), coin)
		c.coinsBySymbol.Store(coin.GetSymbol(), coin)
	}
}
