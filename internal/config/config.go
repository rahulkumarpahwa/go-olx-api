package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT         string
	READTIMEOUT  time.Duration
	WRITETIMEOUT time.Duration
	IDLETIMEOUT  time.Duration
	APP_STATE    string
}

func MustLoad() (*Config, error) {
	godotenv.Load() // we will not check for the env but for the individual variable to exist and if not we will panic.

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is required.")
	}

	app_state := os.Getenv("APP_STATE")
	if app_state == "" {
		panic("APP_STATE is required.")
	}

	readTimeout, err := durationFromEnv("READTIMEOUT")
	if err != nil {
		return nil, err
	}

	writeTimeout, err := durationFromEnv("WRITETIMEOUT")
	if err != nil {
		return nil, err
	}

	idleTimeout, err := durationFromEnv("IDLETIMEOUT")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		PORT:         port,
		READTIMEOUT:  readTimeout,
		WRITETIMEOUT: writeTimeout,
		IDLETIMEOUT:  idleTimeout,
		APP_STATE:    app_state,
	}

	return cfg, nil
}

func durationFromEnv(key string) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf(`%v is required.`, key))
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}

	return duration, nil
}
