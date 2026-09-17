package repositories

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/rahulkumarpahwa/go-olx-api/internal/config"
	"github.com/rahulkumarpahwa/go-olx-api/internal/dto/listings"
)

type ListingStorage interface {
	GetListings(ctx context.Context, requestId string) ([]listings.ListingResponse, error)
	DeleteListingById(ctx context.Context, requestId string, id string) error
}

type ListingRepositories struct {
	Config *config.Config
	DB     *sql.DB
	Logger *slog.Logger
}

func NewListingRepository(cfg *config.Config, db *sql.DB, logger *slog.Logger) *ListingRepositories {
	return &ListingRepositories{
		Config: cfg,
		DB:     db,
		Logger: logger,
	}
}

// todo : add pagination and Limit
func (r *ListingRepositories) GetListings(ctx context.Context, requestId string) ([]listings.ListingResponse, error) {
	const query = "SELECT id, title, description, price, status, city, user_id, category_id, created_at, updated_at FROM listings"

	// const query = "SELECT id, title, description, price, status, city, user_id, category_id, created_at,0 updated_at, pg_sleep(20) FROM listings"

	rows, err := r.DB.QueryContext(ctx, query)

	if err != nil {
		if err == sql.ErrNoRows {
			r.Logger.Error("failed to find rows", "request_id", requestId, "err", err)
			return nil, err
		}

		r.Logger.Error("listings query error", "request_id", requestId, "err", err)
		return nil, err
	}

	defer rows.Close()

	var lstngs []listings.ListingResponse

	for rows.Next() {
		var l listings.ListingResponse
		err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.Status, &l.City, &l.UserID, &l.CategoryID, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			r.Logger.Error("rows scan error", "request_id", requestId, "err", err)
			return nil, err
		}
		r.Logger.Info("listings fetched", "total", len(lstngs))
		lstngs = append(lstngs, l)
	}

	err = rows.Err()
	if err != nil {
		r.Logger.Error("scanned user row error", "request_id", requestId, "err", err)
		return nil, err
	}

	return lstngs, nil
}

func (r *ListingRepositories) DeleteListingById(ctx context.Context, requestId string, id string) error {
	const query = "DELETE FROM listings WHERE id=$1"
	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		r.Logger.Error("listing delete failed", "listing_id", id, "request_id", requestId, "err", err)
		return err
	}
	return nil
}
