package handlers

import (
	"net/http"

	"github.com/rahulkumarpahwa/go-olx-api/internal/httpx"
)

func (h *HealthHandlers) Health(w http.ResponseWriter, r *http.Request) {
	health := h.HealthServices.Health()
	httpx.Write(w, http.StatusOK, health)
}
