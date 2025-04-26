package util

import (
	"log"
	"os"
	"strings"

	"github.com/getsentry/sentry-go"
)

func StartSentry(dsn string) error {
	releaseVersion := os.Getenv("SOURCE_COMMIT")
	log.Printf("Release version: %s", releaseVersion)
	return sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Debug:            false,
		AttachStacktrace: true,
		SampleRate:       1.0,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		SendDefaultPII:   true,
		Release:          releaseVersion,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if strings.Contains(event.Message, "incorrect username/password") {
				return nil
			}
			return event
		},
	})
}
