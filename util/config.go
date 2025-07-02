package util

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server *Server

	DB *Database

	// The duration of the sleep between the requests to the LibreView API.
	FetchSleepDuration time.Duration

	// Sentry configutation URL.
	SentryDSN string
}

type Server struct {
	// The hostname of this application.
	Host string

	// The port to listen on.
	Port uint16

	// Medsenger Agent secret key.
	MedsengerAgentKey string

	// Sets server to debug mode.
	Debug bool
}

type Database struct {
	User string

	Password string

	Dbname string

	Host string
}

func LoadConfigFromEnv() *Config {
	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		panic(err)
	}
	fetchSleepDuration, err := strconv.Atoi(os.Getenv("FETCH_SLEEP_DURATION_MIN"))
	if err != nil {
		panic(err)
	}
	return &Config{
		Server: &Server{
			Host:              os.Getenv("SERVER_HOST"),
			Port:              uint16(serverPort),
			MedsengerAgentKey: os.Getenv("FREESTYLE_LIBRE_KEY"),
			Debug:             os.Getenv("DEBUG") == "true",
		},
		DB: &Database{
			User:     os.Getenv("DB_LOGIN"),
			Password: os.Getenv("DB_PASSWORD"),
			Dbname:   os.Getenv("DB_DATABASE"),
			Host:     os.Getenv("DB_HOST"),
		},
		FetchSleepDuration: time.Duration(fetchSleepDuration) * time.Minute,
		SentryDSN:          os.Getenv("SENTRY_DSN"),
	}
}
