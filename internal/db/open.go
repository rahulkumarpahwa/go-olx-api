package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rahulkumarpahwa/go-olx-api/internal/config"
)

func Open(cfg *config.Config) (*sql.DB, error) {

	db, err := sql.Open("pgx", cfg.DATABASE_URL)
	if err != nil {
		return nil, fmt.Errorf("DB Connection Error: %w", err)
	}

	maxOpen := cfg.MAXOPENCONN
	if maxOpen == 0 {
		maxOpen = 25
	}

	maxIdle := cfg.MAXIDLECONN
	if maxIdle == 0 {
		maxIdle = 25
	}

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(cfg.CONNMAXLIFETIME)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("DB Ping Error: %w", err)
	}

	return db, nil
}
