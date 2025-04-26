package main

import (
	"context"
	"log"
	"time"

	"github.com/TikhonP/maigo"
	"github.com/getsentry/sentry-go"
	"github.com/tikhonp/medsenger-freestylelibre-bot/config"
	"github.com/tikhonp/medsenger-freestylelibre-bot/db"
	"github.com/tikhonp/medsenger-freestylelibre-bot/util"
)

func task(mc *maigo.Client) error {
	lcs, err := db.GetActiveLibreClientToFetch()
	if err != nil {
		return err
	}
	for _, lc := range lcs {
		err := lc.FetchData(mc)
		if err != nil {
			return err
		}
	}
	return nil
}

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
	client := maigo.Init(cfg.Server.MedsengerAgentKey)

	sleepDuration := cfg.FetchSleepDuration.GoDuration()
	ticker := time.NewTicker(sleepDuration)
	for {
		err := task(client)
		if err != nil {
			sentry.CaptureException(err)
			log.Println("Error:", err)
		}
		log.Println("Task completed. Sleeping for", sleepDuration)
		<-ticker.C
	}
}
