package api

import (
	"fmt"
	"github.com/MinterTeam/swap-router/app/config"
	"github.com/MinterTeam/swap-router/app/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"net/http"
)

type Api struct {
	swapService *services.Swap
	coinService *services.Coin
}

func NewApi(cfg config.ApiConfig, ss *services.Swap, cc *services.Coin) *Api {
	api := &Api{
		swapService: ss,
		coinService: cc,
	}

	router := SetupRouter(api)
	if err := router.Run(fmt.Sprintf(":%d", cfg.ServerPort)); err != nil {
		log.WithError(err).Fatal("failed to run api server")
	}

	return api
}

// SetupRouter configure gin router
func SetupRouter(api *Api) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(cors.Default())
	//router.Use(gin.Logger())
	router.Use(gin.ErrorLogger()) // print all errors
	router.Use(apiRecovery)       // returns 500 on any code panics

	// metrics
	p := ginprometheus.NewPrometheus("gin")
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string { return "" } // do not save stats for all routes
	p.Use(router)

	// Default handler 404
	router.NoRoute(func(c *gin.Context) {
		errorResponse(http.StatusNotFound, "Resource not found.", c)
	})

	// Create routing
	router.GET("/status", api.Status)
	router.GET("/api/v1/pools/:coin0/:coin1/route", api.FindSwapRoute)

	return router
}

// Send 500 status and JSON response
func apiRecovery(c *gin.Context) {
	defer func(c *gin.Context) {
		if rec := recover(); rec != nil {
			log.WithField("err", rec).Error("API error")
			errorResponse(http.StatusInternalServerError, "Internal server error", c)
		}
	}(c)

	c.Next()
}
