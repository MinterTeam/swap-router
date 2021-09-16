package api

import (
	"errors"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/swap-router/app/resources"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (api *Api) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (api *Api) FindSwapRoute(c *gin.Context) {
	var req FindSwapPoolRouteRequest
	if err := c.ShouldBindUri(&req); err != nil {
		validationErrorResponse(err, c)
		return
	}

	var reqQuery FindSwapPoolRouteRequestQuery
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		validationErrorResponse(err, c)
		return
	}

	if req.Coin0 == req.Coin1 {
		errorResponse(http.StatusNotFound, "Route path not exists.", c)
		return
	}

	fromCoinId, toCoinId, err := api.getCoinsFromRequest(req)
	if err != nil {
		errorResponse(http.StatusNotFound, err.Error(), c)
		return
	}

	trade, err := api.swapService.FindRoute(fromCoinId, toCoinId, reqQuery.GetTradeType(), reqQuery.GetAmount())
	if err != nil {
		errorResponse(http.StatusNotFound, "Route path not exists.", c)
		return
	}

	path := make([]models.Coin, len(trade.Route.Path))
	for i, t := range trade.Route.Path {
		coin := api.coinService.GetCoinById(t.CoinID)
		path[i] = coin
	}

	c.JSON(http.StatusOK, new(resources.Route).Transform(path, trade))
}

func (api *Api) getCoinsFromRequest(req FindSwapPoolRouteRequest) (fromCoinId, toCoinId uint64, err error) {
	if id, err := strconv.ParseUint(req.Coin0, 10, 64); err == nil {
		fromCoinId = id
	} else {
		if fromCoinId, err = api.coinRepository.FindIdBySymbol(req.Coin0); err != nil {
			return fromCoinId, toCoinId, errors.New("path not found")
		}
	}

	if id, err := strconv.ParseUint(req.Coin1, 10, 64); err == nil {
		toCoinId = id
	} else {
		if toCoinId, err = api.coinRepository.FindIdBySymbol(req.Coin1); err != nil {
			return fromCoinId, toCoinId, errors.New("path not found")
		}
	}

	return fromCoinId, toCoinId, nil
}
