package cfg

import (
	"log"
	"os"
	"strconv"
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

func Load() (*Config, error) {

	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	cfg := &Config{
		PORT:         os.Getenv("PORT"),
		READTIMEOUT:  timeConversion("READTIMEOUT"),
		WRITETIMEOUT: timeConversion("WRITETIMEOUT"),
		IDLETIMEOUT:  timeConversion("IDLETIMEOUT"),
	}
	return cfg, nil
}

func timeConversion(param string) time.Duration {
	readTimeOut, err := strconv.Atoi(os.Getenv(param))
	if(err!=nil){
		log.Fatalf("Parsing Int Error: %v", err)
	}
	return time.Second * time.Duration(readTimeOut)
}
