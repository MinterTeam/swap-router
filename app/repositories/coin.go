package repositories

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v10"
	"strconv"
)

type Coin struct {
	db *pg.DB
}

func NewCoinRepository(db *pg.DB) *Coin {
	return &Coin{db}
}

func (c *Coin) GetBySymbolAndVersion(symbol string, version *uint64) []models.Coin {
	var coins []models.Coin

	query := c.db.Model(&coins).
		Where("symbol = ?", symbol).
		Where("deleted_at IS NULL").
		OrderExpr(`case when "coin"."id" = 0 then 0 else 1 end`).
		Order("reserve DESC")

	if version != nil {
		query.Where("version = ?", version)
	}

	err := query.Select()
	helpers.CheckErr(err)

	return coins
}

func (c *Coin) FindIdBySymbol(symbol string) (uint64, error) {
	if id, err := strconv.ParseUint(symbol, 10, 64); err != nil {
		symbol, version := helpers.GetSymbolAndDefaultVersionFromStr(symbol)
		coins := c.GetBySymbolAndVersion(symbol, &version)
		if len(coins) == 0 {
			return 0, pg.ErrNoRows
		}

		return uint64(coins[0].ID), nil
	} else {
		return id, nil
	}
}

func (c *Coin) GetAll() (coins []models.Coin, err error) {
	err = c.db.Model(&coins).
		Where("deleted_at IS NULL").
		OrderExpr(`case when "coin"."id" = 0 then 0 else 1 end`).
		Order("reserve DESC").
		Select()
	return
}
