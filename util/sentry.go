package util

import (
	"log"

	"github.com/getsentry/sentry-go"
)

func SentryInit(dsn string) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Debug:            false,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		SendDefaultPII:   true,
	}); err != nil {
		log.Printf("Sentry initialization failed: %v", err)
	}
}