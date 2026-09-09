package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rahulkumarpahwa/go-olx-api/internal/httpx"
	"github.com/rahulkumarpahwa/go-olx-api/internal/middleware"
	"github.com/rahulkumarpahwa/go-olx-api/internal/types"
)


func (lh *Handlers) GetListings(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	const query = "SELECT id, title, description, price, status, city, user_id, category_id, created_at, updated_at FROM listings"

	// const query = "SELECT id, title, description, price, status, city, user_id, category_id, created_at,0 updated_at, pg_sleep(20) FROM listings"

	rows, err := lh.DB.QueryContext(r.Context(), query)
	if err != nil {
		if err == sql.ErrNoRows {
			lh.Logger.Error("failed to find rows", "request_id", requestId, "err", err)
			httpx.Error(w, http.StatusNoContent, "something went wrong", httpx.NotFound)
			return
		}

		lh.Logger.Error("listings query error", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.InternalError)
		return
	}

	defer rows.Close()

	var listings []types.Listings

	for rows.Next() {
		var l types.Listings
		err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.Status, &l.City, &l.UserID, &l.CategoryID, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			lh.Logger.Error("rows scan error", "request_id", requestId, "err", err)
			httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.InternalError)
			return
		}
		lh.Logger.Info("listings fetched", "total", len(listings))
		listings = append(listings, l)
	}

	err = rows.Err()
	if err != nil {
		lh.Logger.Error("scanned rows error", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.InternalError)
		return
	}

	httpx.Write(w, http.StatusOK, listings)
}

func (lh *Handlers) DeleteListing(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	if id == "" {
		lh.Logger.Error("missing delete listing id", "listing_id", id, "request_id", requestId)
		httpx.Error(w, http.StatusBadRequest, "missing delete listing id", httpx.InvalidId)
		return
	}

	const query = "DELETE FROM listings WHERE id=$1"

	_, err := lh.DB.ExecContext(r.Context(), query, id)
	if err != nil {
		lh.Logger.Error("listing delete failed", "listing_id", id, "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong	", httpx.InternalError)
		return
	}
	// when still we get the error in delete still we will return the true to improve the security of the app.

	httpx.Write(w, http.StatusNoContent, "listing deleted successfully")
}

func (lh *Handlers) CreateListing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB max
	defer r.Body.Close()

	var body CreateListingRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		lh.Logger.Error("failed to decode", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "inavlid body", httpx.MalformedJSON)
		return // must
	}

	// basic validation
	if body.Title == "" || body.City == "" || body.Price <= 0 || body.UserID == uuid.Nil || body.CategoryID == uuid.Nil {
		lh.Logger.Error("missing or invalid fields", "request_id", requestId)
		httpx.Error(w, http.StatusBadRequest, "missing or invalid fields", httpx.ValidationFailed)
		return
	}

	const query = "INSERT INTO listings(title, description, price, city, user_id, category_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, status, created_at, updated_at"

	var (
		id         uuid.UUID
		status     types.ListingStatus
		created_at time.Time
		updated_at time.Time
	)

	err := lh.DB.QueryRowContext(r.Context(), query, &body.Title, &body.Description, &body.Price, &body.City, &body.UserID, &body.CategoryID).Scan(&id, &status, &created_at, &updated_at)
	if err != nil {
		lh.Logger.Error("failed to scan row", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "failed to scan row", httpx.BadRequest)
		return
	}

	body.ID = id
	body.Status = status
	body.CreatedAt = created_at
	body.UpdatedAt = &updated_at

	lh.Logger.Info("listing created successfully", "request_id", requestId, "listing_id", id)

	httpx.Write(w, http.StatusCreated, body)
}
