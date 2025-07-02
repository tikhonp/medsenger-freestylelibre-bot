package main

import (
	"errors"
	"log"
	"slices"
	"time"

	"github.com/TikhonP/maigo"
	"github.com/getsentry/sentry-go"
	"github.com/tikhonp/medsenger-freestylelibre-bot/db"
	"github.com/tikhonp/medsenger-freestylelibre-bot/util"
	libreclient "github.com/tikhonp/medsenger-freestylelibre-bot/util/libre_client"
)

var sentryExcludedErrors = []error{
	libreclient.ErrIncorrectUsernameOrPassword,
	db.ErrLibreAccountConnectionsIsEmpty,
}

func processTaskError(err error) {
	if err == nil {
		return
	}
	log.Println("Error:", err)
	errIsRestrictedToSendToSentry := slices.ContainsFunc(sentryExcludedErrors, func(rErr error) bool { return errors.Is(err, rErr) })
	if !errIsRestrictedToSendToSentry {
		sentry.CaptureException(err)
	}
}

func task(mc *maigo.Client) {
	lcs, err := db.GetActiveLibreClientToFetch()
	if err != nil {
		sentry.CaptureException(err)
		log.Println("Error:", err)
		return
	}
	for _, lc := range lcs {
		err := lc.FetchData(mc)
		processTaskError(err)
	}
}

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
	client := maigo.Init(cfg.Server.MedsengerAgentKey)
	ticker := time.NewTicker(cfg.FetchSleepDuration)
	for {
		task(client)
		log.Println("Task completed. Sleeping for", cfg.FetchSleepDuration)
		<-ticker.C
	}
}
