package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type Config struct {
	DbConfig      DbConfig
	ApiConfig     ApiConfig
	WsConfig      WsConfig
	WorkersConfig WorkersConfig
}

type DbConfig struct {
	Host        string
	Name        string
	User        string
	Password    string
	PoolSize    int
	SslRequired bool
}

type ApiConfig struct {
	ServerPort uint64
}

type WsConfig struct {
	Server string
}

type WorkersConfig struct {
	FindRouteWorkersCount int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Warn(".env file not found")
	}

	serverPort, _ := strconv.ParseUint(os.Getenv("SERVER_PORT"), 10, 64)
	poolSize, _ := strconv.Atoi(os.Getenv("DB_POOL_SIZE"))
	sslRequired, _ := strconv.ParseBool(os.Getenv("DB_SSL_REQUIRED"))
	findRouteWorkersCount, _ := strconv.Atoi(os.Getenv("FIND_ROUTE_WORKERS_COUNT"))

	return &Config{
		DbConfig: DbConfig{
			Host:        os.Getenv("DB_HOST"),
			Name:        os.Getenv("DB_NAME"),
			User:        os.Getenv("DB_USER"),
			Password:    os.Getenv("DB_PASSWORD"),
			SslRequired: sslRequired,
			PoolSize:    poolSize,
		},
		ApiConfig: ApiConfig{
			ServerPort: serverPort,
		},
		WsConfig: WsConfig{
			Server: os.Getenv("WS_SERVER"),
		},
		WorkersConfig: WorkersConfig{
			FindRouteWorkersCount: findRouteWorkersCount,
		},
	}
}
