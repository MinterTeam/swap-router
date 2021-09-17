package services

import (
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/swap-router/app/repositories"
	log "github.com/sirupsen/logrus"
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
	log.Debugf("coin: received new block %d", b.ID)
	c.fillCoinsMap()
}

func (c *Coin) fillCoinsMap() {
	wg := &sync.WaitGroup{}
	coins, _ := c.repository.GetAll()
	for _, coin := range coins {
		wg.Add(1)
		go func(wg *sync.WaitGroup, coin models.Coin) {
			defer wg.Done()
			c.coins.Store(uint64(coin.ID), coin)
		}(wg, coin)
	}
	wg.Wait()
}
