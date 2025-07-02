package main

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	libre "github.com/tikhonp/medsenger-freestylelibre-bot"
	"github.com/tikhonp/medsenger-freestylelibre-bot/db"
	"github.com/tikhonp/medsenger-freestylelibre-bot/util"
)

func main() {
	cfg := util.LoadConfigFromEnv()
	if !cfg.Server.Debug {
		err := util.StartSentry(cfg.SentryDSN)
		if err != nil {
			log.Fatalln(err)
		}
		defer sentry.Flush(2 * time.Second)
	}
	db.MustConnect(cfg.DB)
	libre.NewServer(cfg.Server).Listen()
}
