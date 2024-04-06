package main

import (
	"context"

	"github.com/TikhonP/medsenger-freestylelibre-bot/config"
	"github.com/TikhonP/medsenger-freestylelibre-bot/db"
    libre "github.com/TikhonP/medsenger-freestylelibre-bot"
)

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "pkl/local/config.pkl")
	if err != nil {
		panic(err)
	}
	db.Connect(cfg.Db)
    libre.Serve(cfg.Server)
}
