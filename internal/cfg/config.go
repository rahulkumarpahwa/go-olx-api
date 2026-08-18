package cfg

import "time"

type CONFIG struct {
	PORT         string
	READTIMEOUT  time.Duration
	WRITETIMEOUT time.Duration
	IDLETIMEOUT  time.Duration
}

func Config() CONFIG {
	cfg := CONFIG{
		PORT:         ":8080",
		READTIMEOUT:  time.Second * 20,
		WRITETIMEOUT: time.Second * 20,
		IDLETIMEOUT:  time.Second * 10,
	}
	return cfg
}
