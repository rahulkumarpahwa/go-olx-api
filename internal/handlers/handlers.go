package handlers

import (
	"database/sql"
	"log/slog"

	"github.com/rahulkumarpahwa/go-olx-api/internal/config"
)

type Handlers struct {
	Config *config.Config
	DB     *sql.DB
	Logger *slog.Logger
}

func NewHanlders(cfg *config.Config, db *sql.DB, logger *slog.Logger) *Handlers {
	return &Handlers{
		Config: cfg,
		DB:     db,
		Logger: logger,
	}
}



func NewUserHanlders() *Handlers {
	return &Handlers{
	}
}
