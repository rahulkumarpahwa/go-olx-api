package handlers

import (
	"log/slog"

	"github.com/rahulkumarpahwa/go-olx-api/internal/config"
	"github.com/rahulkumarpahwa/go-olx-api/internal/services"
)

type HealthHandlers struct {
	HealthServices *services.HealthServices
	Logger         *slog.Logger
}

type UserHandlers struct {
	UserServices *services.UserServices
	Config       *config.Config
	Logger       *slog.Logger
}

type ListingHandlers struct {
	ListingServices *services.ListingServices
	Logger          *slog.Logger
}

type Handlers struct {
	Config *config.Config
	Logger *slog.Logger
}

func NewHandler(cfg *config.Config, logger *slog.Logger) *Handlers {
	return &Handlers{
		Config: cfg,
		Logger: logger,
	}
}

func (h *Handlers) UserHandlers(userServices *services.UserServices) *UserHandlers {
	return &UserHandlers{
		UserServices: userServices,
		Logger:       h.Logger,
	}
}

func (h *Handlers) ListingHanlders(listingServices *services.ListingServices) *ListingHandlers {
	return &ListingHandlers{
		ListingServices: listingServices,
		Logger:          h.Logger,
	}
}

func (h *Handlers) HealthHandlers(healthServices *services.HealthServices) *HealthHandlers {
	return &HealthHandlers{
		HealthServices: healthServices,
		Logger:         h.Logger,
	}
}
