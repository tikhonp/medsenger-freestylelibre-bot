package main

import (
	"context"

	libre "github.com/TikhonP/medsenger-freestylelibre-bot"
	"github.com/TikhonP/medsenger-freestylelibre-bot/config"
	"github.com/TikhonP/medsenger-freestylelibre-bot/db"
)

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "pkl/local/config.pkl")
	if err != nil {
		panic(err)
	}
	db.Connect(cfg.Db)
	libre.Serve(cfg.Server)
}
