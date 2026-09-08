package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT               string
	READTIMEOUT        time.Duration
	WRITETIMEOUT       time.Duration
	IDLETIMEOUT        time.Duration
	APP_STATE          string
	DATABASE_URL       string
	MAXOPENCONN        int
	MAXIDLECONN        int
	CONNMAXLIFETIME    time.Duration
	MIGRATIONFILESPATH string
	TIMEZONE           string
}

func MustLoad() (*Config, error) {
	godotenv.Load() // we will not check for the env but for the individual variable to exist and if not we will panic.

	port, ok := os.LookupEnv("PORT")
	if !ok {
		panic("PORT is required.")
	}

	app_state, ok := os.LookupEnv("APP_STATE")
	if !ok {
		panic("APP_STATE is required.")
	}

	database_url, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		panic("DATABASE_URL is required.")
	}

	max_open_conn, err := convertToInt("MAXOPENCONN")
	if err != nil {
		return nil, err
	}

	max_idle_conn, err := convertToInt("MAXIDLECONN")
	if err != nil {
		return nil, err
	}

	conn_max_lifetime, err := durationFromEnv("CONNMAXLIFETIME")
	if err != nil {
		return nil, err
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

	migration_file_path, ok := os.LookupEnv("MIGRATIONFILESPATH")
	if !ok {
		return nil, fmt.Errorf("MIGRATIONFILESPATH is required.")
	}

	cfg := &Config{
		PORT:               port,
		READTIMEOUT:        readTimeout,
		WRITETIMEOUT:       writeTimeout,
		IDLETIMEOUT:        idleTimeout,
		APP_STATE:          app_state,
		DATABASE_URL:       database_url,
		MAXOPENCONN:        max_open_conn,
		MAXIDLECONN:        max_idle_conn,
		CONNMAXLIFETIME:    conn_max_lifetime,
		MIGRATIONFILESPATH: migration_file_path,
		TIMEZONE:           os.Getenv("TIMEZONE"),
	}

	return cfg, nil
}

func convertToInt(key string) (int, error) {
	value := os.Getenv(key)

	if value == "" {
		return 0, fmt.Errorf(`%s is required.`, key)
	}

	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return num, nil
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

func (c *Config) IsProduction() bool {
	return c.APP_STATE == "Prod"
}

func (c *Config) IsDev() bool {
	return c.APP_STATE == "Dev"
}
