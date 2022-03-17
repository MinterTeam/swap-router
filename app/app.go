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
	// Init Logger
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})

	log.Debugf("num of cpu: %d", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg := config.Load()
	db := database.Connect(cfg.DbConfig)

	poolRepository := repositories.NewPoolRepository(db)
	poolService := services.NewPoolService(poolRepository)
	log.Debug("pool service started")
	swapService := services.NewSwapService(cfg.WorkersConfig, poolService)
	log.Debug("swap service started")
	coinRepository := repositories.NewCoinRepository(db)
	coinService := services.NewCoinService(coinRepository)
	log.Debug("coin service started")

	go func() {
		wsClient := ws.NewWebSocketClient(cfg.WsConfig.Server)
		wsSub := wsClient.CreateSubscription("blocks")
		wsClient.Subscribe(wsSub)
		blocksListener := ws.NewBlocksChannelHandler()
		blocksListener.AddSubscriber(poolService)
		blocksListener.AddSubscriber(coinService)
		wsSub.OnPublish(blocksListener)
		defer wsClient.Close()

		select {}
	}()

	api.NewApi(cfg.ApiConfig, swapService, coinService)
}
