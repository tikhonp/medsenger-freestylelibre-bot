package util

import (
	"log"
	"os"

	"github.com/getsentry/sentry-go"
)

func StartSentry(dsn string, releaseVersionFile string) {
	releaseVersion, err := os.ReadFile(releaseVersionFile)
	if err != nil {
		log.Printf("Failed to read release version file: %v", err)
		releaseVersion = []byte("unknown")
	}
	log.Printf("Release version: %s", string(releaseVersion))
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Debug:            false,
		AttachStacktrace: true,
		SampleRate:       1.0,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		SendDefaultPII:   true,
		Release:          string(releaseVersion),
	}); err != nil {
		log.Printf("Sentry initialization failed: %v", err)
	}
}
