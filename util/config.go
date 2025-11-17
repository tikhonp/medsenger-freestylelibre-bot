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

	LibreClientLLUVersion string
}

type Server struct {
	// The port to listen on.
	Port uint16

	// Medsenger Agent secret key.
	MedsengerAgentKey string

	// Sets server to debug mode.
	Debug bool
}

type Database struct {
	Host string

	Port int

	User string

	Password string

	Dbname string
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
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic(err)
	}
	return &Config{
		Server: &Server{
			Port:              uint16(serverPort),
			MedsengerAgentKey: os.Getenv("FREESTYLE_LIBRE_KEY"),
			Debug:             os.Getenv("DEBUG") == "true",
		},
		DB: &Database{
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			User:     os.Getenv("DB_LOGIN"),
			Password: os.Getenv("DB_PASSWORD"),
			Dbname:   os.Getenv("DB_DATABASE"),
		},
		FetchSleepDuration:    time.Duration(fetchSleepDuration) * time.Minute,
		SentryDSN:             os.Getenv("SENTRY_DSN"),
		LibreClientLLUVersion: os.Getenv("LIBRE_CLIENT_LLU_VERSION"),
	}
}
