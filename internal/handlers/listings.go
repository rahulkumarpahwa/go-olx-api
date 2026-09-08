package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rahulkumarpahwa/go-olx-api/internal/httpx"
	"github.com/rahulkumarpahwa/go-olx-api/internal/middleware"
	"github.com/rahulkumarpahwa/go-olx-api/internal/types"
)

type RequestBody struct {
	ID          uuid.UUID           `json:"-"`
	Title       string              `json:"title"`
	Description *string             `json:"description"`
	Price       int64               `json:"price"`
	City        string              `json:"city"`
	Status      types.ListingStatus `json:"-"`
	UserID      uuid.UUID           `json:"user_id"`
	CategoryID  uuid.UUID           `json:"category_id"`
	CreatedAt   time.Time           `json:"-"`
	UpdatedAt   *time.Time          `json:"-"`
}

func (h *Handlers) GetListings(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	// const query = "SELECT id, title, description, price, status, city, user_id, category_id, created_at, updated_at FROM listings"

	const query = "SELECT id, title, description, price, status, city, user_id, category_id, created_at,0 updated_at, pg_sleep(20) FROM listings"

	rows, err := h.DB.QueryContext(r.Context(), query)
	if err != nil {
		h.Logger.Error("listings query error", "request_id", requestId, "err", err)

		if err == sql.ErrNoRows {
			httpx.Error(w, http.StatusNoContent, "h.DB.QueryContext: No Rows", httpx.NotFound)
			return
		}

		httpx.Error(w, http.StatusInternalServerError, "h.DB.QueryContext: "+err.Error(), httpx.InternalError)
		return
	}

	defer rows.Close()

	var listings []types.Listings

	for rows.Next() {
		var l types.Listings
		err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.Status, &l.City, &l.UserID, &l.CategoryID, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			h.Logger.Error("rows scan error", "request_id", requestId, "err", err)
			httpx.Error(w, http.StatusInternalServerError, "rows.scan: "+err.Error(), httpx.InternalError)
			return
		}
		h.Logger.Info("listings fetched", "total", len(listings))
		listings = append(listings, l)
	}

	err = rows.Err()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "rows.Err(): "+err.Error(), httpx.InternalError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]types.Listings{"listings": listings})
}

func (h *Handlers) DeleteListing(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	if id == "" {
		httpx.Error(w, http.StatusBadRequest, "Missing Delete Listing ID", httpx.InvalidId)
		return
	}

	const query = "DELETE FROM listings WHERE id=$1"

	_, err := h.DB.ExecContext(r.Context(), query, id)
	if err != nil {
		h.Logger.Error("delete failed", "listing_id", id, "request_id", requestId, "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	// when still we get the error in delete still we will return the true to improve the security of the app.

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]any{"status": "Ok"})
}

func (h *Handlers) CreateListing(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB max

	var body RequestBody
	data, err := io.ReadAll(r.Body)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Error reading request body", httpx.MalformedJSON)
		return
	}
	if err := json.Unmarshal(data, &body); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// basic validation
	if body.Title == "" || body.City == "" || body.Price <= 0 || body.UserID == uuid.Nil || body.CategoryID == uuid.Nil {
		httpx.Error(w, http.StatusBadRequest, "Missing or invalid fields", httpx.ValidationFailed)
		return
	}

	const query = "INSERT INTO listings(title, description, price, city, user_id, category_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, status, created_at , updated_at"

	var id uuid.UUID
	var status types.ListingStatus
	var created_at time.Time
	var updated_at time.Time
	err = h.DB.QueryRowContext(r.Context(), query, &body.Title, &body.Description, &body.Price, &body.City, &body.UserID, &body.CategoryID, &body.CreatedAt).Scan(&id, &status, &created_at, &updated_at)
	if err != nil {
		http.Error(w, "h.DB.QueryRowContext: "+err.Error(), http.StatusInternalServerError)
		return
	}

	body.ID = id
	body.Status = status
	body.CreatedAt = created_at
	body.UpdatedAt = &updated_at

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"status": "Ok", "listing": body})
}
