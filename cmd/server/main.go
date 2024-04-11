package main

import (
	"context"
	"log"

	libre "github.com/TikhonP/medsenger-freestylelibre-bot"
	"github.com/TikhonP/medsenger-freestylelibre-bot/config"
	"github.com/TikhonP/medsenger-freestylelibre-bot/db"
	"github.com/getsentry/sentry-go"
)

func sentryInit(dsn string) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		TracesSampleRate: 1.0,
		Debug:            false,
	}); err != nil {
		log.Printf("Sentry initialization failed: %v", err)
	}
}

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "pkl/local/config.pkl")
	if err != nil {
		panic(err)
	}
	if !cfg.Server.Debug {
		sentryInit(cfg.SentryDSN)
	}
	db.Connect(cfg.Db)
	libre.Serve(cfg.Server)
}
