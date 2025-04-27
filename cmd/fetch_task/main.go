package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/TikhonP/maigo"
	"github.com/getsentry/sentry-go"
	"github.com/tikhonp/medsenger-freestylelibre-bot/config"
	"github.com/tikhonp/medsenger-freestylelibre-bot/db"
	"github.com/tikhonp/medsenger-freestylelibre-bot/util"
	libreclient "github.com/tikhonp/medsenger-freestylelibre-bot/util/libre_client"
)

func task(mc *maigo.Client) {
	lcs, err := db.GetActiveLibreClientToFetch()
	if err != nil {
		sentry.CaptureException(err)
		log.Println("Error:", err)
		return
	}
	for _, lc := range lcs {
		err := lc.FetchData(mc)
		if err != nil && !errors.Is(err, libreclient.ErrIncorrectUsernameOrPassword) {
			sentry.CaptureException(err)
			log.Println("Error:", err)
		}
	}
}

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "config.pkl")
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
		task(client)
		log.Println("Task completed. Sleeping for", sleepDuration)
		<-ticker.C
	}
}
