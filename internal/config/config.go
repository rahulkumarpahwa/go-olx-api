package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT            string
	READTIMEOUT     time.Duration
	WRITETIMEOUT    time.Duration
	IDLETIMEOUT     time.Duration
	APP_STATE       string
	DATABASE_URL    string
	MAXOPENCONN     int
	MAXIDLECONN     int
	CONNMAXLIFETIME time.Duration
	MIGRATIONFILESPATH string
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

	database_url := os.Getenv("DATABASE_URL")
	if database_url == "" {
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

	migration_file_path := os.Getenv("MIGRATIONFILESPATH");
	if(migration_file_path == ""){
		return nil, fmt.Errorf("MIGRATIONFILESPATH is required.")
	}

	cfg := &Config{
		PORT:            port,
		READTIMEOUT:     readTimeout,
		WRITETIMEOUT:    writeTimeout,
		IDLETIMEOUT:     idleTimeout,
		APP_STATE:       app_state,
		DATABASE_URL:    database_url,
		MAXOPENCONN:     max_open_conn,
		MAXIDLECONN:     max_idle_conn,
		CONNMAXLIFETIME: conn_max_lifetime,
		MIGRATIONFILESPATH : migration_file_path,
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
