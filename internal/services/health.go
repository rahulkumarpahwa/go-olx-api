package services

import (
	"time"

	"github.com/rahulkumarpahwa/go-olx-api/internal/config"
)

type HealthServices struct {
	Cfg *config.Config
}

func NewHealthServices(cfg *config.Config) *HealthServices {
	return &HealthServices{
		Cfg: cfg,
	}
}

func (h *HealthServices) Health() any {
	health := struct {
		Message     string `json:"message"`
		CurrentTime string `json:"current_time"`
	}{
		Message:     "Server is working fine.",
		CurrentTime: time.Now().In(time.FixedZone(h.Cfg.TIMEZONE, 5*60*60)).Format("2006-01-02 15:04:05"),
	}
	return health
}
