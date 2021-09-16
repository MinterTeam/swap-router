package repositories

import (
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v10"
)

type Pool struct {
	db *pg.DB
}

func NewPoolRepository(db *pg.DB) *Pool {
	return &Pool{db}
}

func (p *Pool) GetAll() (pools []models.LiquidityPool, err error) {
	err = p.db.Model(&pools).Select()
	return
}
