package util

import (
	"log"

	"github.com/getsentry/sentry-go"
)

func StartSentry(dsn string) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Debug:            false,
		AttachStacktrace: true,
		SampleRate:       1.0,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		SendDefaultPII:   true,
	}); err != nil {
		log.Printf("Sentry initialization failed: %v", err)
	}
}
