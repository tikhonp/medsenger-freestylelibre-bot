package main

import (
	"context"
	"log"
	"time"

	"github.com/TikhonP/maigo"
	"github.com/TikhonP/medsenger-freestylelibre-bot/config"
	"github.com/TikhonP/medsenger-freestylelibre-bot/db"
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
	db.Connect(cfg.Db)
	client := maigo.Init(cfg.Server.MedsengerAgentKey)

	sleepDuration := cfg.FetchSleepDuration.GoDuration()
	for {
		err := task(client)
		if err != nil {
			panic(err)
		}
		log.Println("Task completed. Sleeping for", sleepDuration)
		time.Sleep(sleepDuration)
	}
}
