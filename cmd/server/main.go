package main

import (
	"context"

	libre "github.com/TikhonP/medsenger-freestylelibre-bot"
	"github.com/TikhonP/medsenger-freestylelibre-bot/config"
	"github.com/TikhonP/medsenger-freestylelibre-bot/db"
	"github.com/TikhonP/medsenger-freestylelibre-bot/util"
)

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "pkl/local/config.pkl")
	if err != nil {
		panic(err)
	}
	if !cfg.Server.Debug {
		util.StartSentry(cfg.SentryDSN)
	}
	db.MustConnect(cfg.Db)
	libre.NewServer(cfg.Server).Listen()
}
