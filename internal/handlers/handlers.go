package handlers

import (
	"log/slog"

	"github.com/rahulkumarpahwa/go-olx-api/internal/services"
)

type Handlers struct {
	HealthServices  *services.HealthServices
	UserServices    *services.UserServices
	ListingServices *services.ListingServices
	Logger          *slog.Logger
}

func NewHandler(logger *slog.Logger) *Handlers {
	return &Handlers{
		Logger: logger,
	}
}

func (h *Handlers) UserHandlers(userServices *services.UserServices) *Handlers {
	return &Handlers{
		UserServices: userServices,
		Logger:       h.Logger,
	}
}

func (h *Handlers) ListingHanlders(listingServices *services.ListingServices) *Handlers {
	return &Handlers{
		ListingServices: listingServices,
		Logger:          h.Logger,
	}
}

func (h *Handlers) HealthHandlers(healthServices *services.HealthServices) *Handlers {
	return &Handlers{
		HealthServices: healthServices,
		Logger:         h.Logger,
	}
}
