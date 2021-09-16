package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type Config struct {
	DbConfig  DbConfig
	ApiConfig ApiConfig
	WsConfig  WsConfig
}

type DbConfig struct {
	Host     string
	Name     string
	User     string
	Password string
	PoolSize int
}

type ApiConfig struct {
	ServerPort uint64
}

type WsConfig struct {
	Server string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Warn(".env file not found")
	}

	serverPort, _ := strconv.ParseUint(os.Getenv("SERVER_PORT"), 10, 64)
	poolSize, _ := strconv.Atoi(os.Getenv("DB_POOL_SIZE"))

	return &Config{
		DbConfig: DbConfig{
			Host:     os.Getenv("DB_HOST"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			PoolSize: poolSize,
		},
		ApiConfig: ApiConfig{
			ServerPort: serverPort,
		},
		WsConfig: WsConfig{
			Server: os.Getenv("WS_SERVER"),
		},
	}
}
