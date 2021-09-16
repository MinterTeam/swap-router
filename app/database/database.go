package database

import (
	"crypto/tls"
	"fmt"
	"github.com/MinterTeam/swap-router/app/config"
	"github.com/go-pg/pg/extra/pgdebug/v10"
	"github.com/go-pg/pg/v10"
	"log"
)

func Connect(cfg config.DbConfig) *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Name,
		Addr:     cfg.Host,
		PoolSize: cfg.PoolSize,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	})

	db.AddQueryHook(pgdebug.NewDebugHook())

	return db
}

func Close(db *pg.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not close connection to database: %s", err))
	}
}
