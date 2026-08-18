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
		READTIMEOUT:  time.Second * 10,
		WRITETIMEOUT: time.Second * 30,
		IDLETIMEOUT:  time.Second * 60,
	}
	return cfg
}
