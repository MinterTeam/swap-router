package main

import (
	"github.com/MinterTeam/swap-router/app/api"
	"github.com/MinterTeam/swap-router/app/config"
	"github.com/MinterTeam/swap-router/app/database"
	"github.com/MinterTeam/swap-router/app/repositories"
	"github.com/MinterTeam/swap-router/app/services"
	"github.com/MinterTeam/swap-router/app/ws"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"time"
)

func main() {
	log.Debugf("num of cpu: %d", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Init Logger
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})

	cfg := config.Load()
	db := database.Connect(cfg.DbConfig)

	poolRepository := repositories.NewPoolRepository(db)
	poolService := services.NewPoolService(poolRepository)
	swapService := services.NewSwapService(cfg.WorkersConfig, poolService)
	coinRepository := repositories.NewCoinRepository(db)
	coinService := services.NewCoinService(coinRepository)

	wsClient := ws.NewWebSocketClient(cfg.WsConfig.Server)
	wsSub := wsClient.CreateSubscription("blocks")
	wsClient.Subscribe(wsSub)
	blocksListener := ws.NewBlocksChannelHandler()
	blocksListener.AddSubscriber(poolService)
	blocksListener.AddSubscriber(coinService)
	wsSub.OnPublish(blocksListener)
	defer wsClient.Close()

	api.NewApi(cfg.ApiConfig, swapService, coinService, coinRepository)
}
