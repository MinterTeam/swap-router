package services

import (
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/swap-router/app/helpers"
	"github.com/MinterTeam/swap-router/app/repositories"
	"github.com/MinterTeam/swap-router/app/swap"
	log "github.com/sirupsen/logrus"
)

type Pool struct {
	repository *repositories.Pool

	pools      []models.LiquidityPool
	tradePairs []swap.Pair
}

func NewPoolService(r *repositories.Pool) *Pool {
	service := &Pool{repository: r}
	service.updatePools()
	service.updateTradePairs()
	return service
}

func (p *Pool) updatePools() {
	pools, err := p.repository.GetAll()
	if err != nil {
		log.Errorf("failed to get pools: %s", err)
		return
	}

	p.pools = pools
}

func (p *Pool) updateTradePairs() {
	pairs := make([]swap.Pair, len(p.pools))
	for i, pool := range p.pools {
		pairs[i] = swap.NewPair(
			swap.NewTokenAmount(swap.NewToken(pool.FirstCoinId), helpers.Str2BigInt(pool.FirstCoinVolume)),
			swap.NewTokenAmount(swap.NewToken(pool.SecondCoinId), helpers.Str2BigInt(pool.SecondCoinVolume)),
		)
	}

	p.tradePairs = pairs
}

func (p *Pool) GetPools() []models.LiquidityPool {
	return p.pools
}

func (p *Pool) GetTradePairs() []swap.Pair {
	return p.tradePairs
}

func (p *Pool) ListenNewBlock(b blocks.Resource) {
	log.Debugf("pool: received new block %d", b.ID)
	p.updatePools()
	p.updateTradePairs()
}
