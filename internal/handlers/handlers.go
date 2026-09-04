package handlers

import (
	"database/sql"
	"log/slog"
)

type Handlers struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewHanlders(db *sql.DB, logger *slog.Logger) *Handlers {
	return &Handlers{
		DB:     db,
		Logger: logger,
	}
}
