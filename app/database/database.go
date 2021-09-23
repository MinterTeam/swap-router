package database

import (
	"crypto/tls"
	"fmt"
	"github.com/MinterTeam/swap-router/app/config"
	"github.com/go-pg/pg/v10"
	"log"
)

func Connect(cfg config.DbConfig) *pg.DB {
	options := &pg.Options{
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Name,
		Addr:     cfg.Host,
		PoolSize: cfg.PoolSize,
	}

	if cfg.SslRequired {
		options.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	db := pg.Connect(options)

	return db
}

func Close(db *pg.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not close connection to database: %s", err))
	}
}
