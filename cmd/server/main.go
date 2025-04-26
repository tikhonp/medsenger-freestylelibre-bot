package main

import (
	"context"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	libre "github.com/tikhonp/medsenger-freestylelibre-bot"
	"github.com/tikhonp/medsenger-freestylelibre-bot/config"
	"github.com/tikhonp/medsenger-freestylelibre-bot/db"
	"github.com/tikhonp/medsenger-freestylelibre-bot/util"
)

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "pkl/local/config.pkl")
	if err != nil {
		panic(err)
	}
	if !cfg.Server.Debug {
		err = util.StartSentry(cfg.SentryDSN)
		if err != nil {
			log.Fatalln(err)
		}
		defer sentry.Flush(2 * time.Second)
	}
	db.MustConnect(cfg.Db)
	libre.NewServer(cfg.Server).Listen()
}
