package services

import (
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/swap-router/app/repositories"
	"sync"
)

type Coin struct {
	repository *repositories.Coin

	coins sync.Map
}

func NewCoinService(r *repositories.Coin) *Coin {
	service := &Coin{repository: r}
	service.fillCoinsMap()
	return service
}

func (c *Coin) GetCoinById(id uint64) models.Coin {
	coin, _ := c.coins.Load(id)
	return coin.(models.Coin)
}

func (c *Coin) ListenNewBlock(b blocks.Resource) {
	c.fillCoinsMap()
}

func (c *Coin) fillCoinsMap() {
	coins, _ := c.repository.GetAll()
	for _, coin := range coins {
		c.coins.Store(uint64(coin.ID), coin)
	}
}
