package handlers

import (
	"net/http"
	"time"

	"github.com/rahulkumarpahwa/go-olx-api/internal/httpx"
)

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	health := struct {
		Message     string `json:"message"`
		CurrentTime string `json:"current_time"`
	}{
		Message:     "Server is working fine.",
		CurrentTime: time.Now().In(time.FixedZone(h.Config.TIMEZONE, 5*60*60)).Format("2006-01-02 15:04:05"),
	}

	httpx.Write(w, http.StatusOK, health)
}
