package main

import (
	"github.com/MinterTeam/swap-router/app/api"
	"github.com/MinterTeam/swap-router/app/config"
	"github.com/MinterTeam/swap-router/app/database"
	"github.com/MinterTeam/swap-router/app/repositories"
	"github.com/MinterTeam/swap-router/app/services"
	"github.com/MinterTeam/swap-router/app/ws"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	cfg := config.Load()
	db := database.Connect(cfg.DbConfig)

	poolRepository := repositories.NewPoolRepository(db)
	poolService := services.NewPoolService(poolRepository)
	swapService := services.NewSwapService(poolService)
	coinRepository := repositories.NewCoinRepository(db)
	coinService := services.NewCoinService(coinRepository)

	wsClient := ws.NewWebSocketClient(cfg.WsConfig.Server)
	defer wsClient.Close()

	blocksListener := ws.NewBlocksChannelHandler()
	blocksListener.AddSubscriber(poolService)

	api.NewApi(cfg.ApiConfig, swapService, coinService, coinRepository)
}
