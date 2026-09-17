package handlers

import (
	"net/http"

	"github.com/rahulkumarpahwa/go-olx-api/internal/httpx"
	"github.com/rahulkumarpahwa/go-olx-api/internal/middleware"
)

func (lh *ListingHandlers) GetAll(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	listings, err := lh.ListingServices.Storage.GetListings(ctx, requestId)

	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.InternalError)
		return
	}

	httpx.Write(w, http.StatusOK, listings)
}

func (lh *ListingHandlers) DeleteById(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	if id == "" {
		lh.Logger.Error("missing delete listing id", "listing_id", id, "request_id", requestId)
		httpx.Error(w, http.StatusBadRequest, "missing delete listing id", httpx.InvalidId)
		return
	}

	err := lh.ListingServices.Storage.DeleteListingById(ctx, requestId, id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.InternalError)
		return
	}
	// when still we get the error in delete still we will return the true to improve the security of the app.

	httpx.Write(w, http.StatusNoContent, "listing deleted successfully")
}

// func (lh *Handlers) Create(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	requestId := middleware.RequestIDFromContext(ctx)

// 	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB max
// 	defer r.Body.Close()

// 	var body dto.CreateListingRequest
// 	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
// 		lh.Logger.Error("failed to decode", "request_id", requestId, "err", err)
// 		httpx.Error(w, http.StatusBadRequest, "inavlid body", httpx.MalformedJSON)
// 		return // must
// 	}

// 	// basic validation
// 	if err := body.Validate(); err != nil {
// 		var verr *dto.ValidationError
// 		errors.As(err, &verr) // this one method opens the error and parse it into the value passed as reference.
// 		lh.Logger.Error("missing or invalid fields", "request_id", requestId, "err", err)
// 		httpx.ValidationError(w, http.StatusBadRequest, err.Error(), httpx.ValidationFailed, verr.Field)
// 		return
// 	}

// 	const query = "INSERT INTO listings(title, description, price, city, user_id, category_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, status, created_at, updated_at"

// 	var (
// 		id         uuid.UUID
// 		status     types.ListingStatus
// 		created_at time.Time
// 		updated_at time.Time
// 	)

// 	err := lh.DB.QueryRowContext(r.Context(), query, &body.Title, &body.Description, &body.Price, &body.City, &body.UserID, &body.CategoryID).Scan(&id, &status, &created_at, &updated_at)
// 	if err != nil {
// 		lh.Logger.Error("failed to scan row", "request_id", requestId, "err", err)
// 		httpx.Error(w, http.StatusBadRequest, "failed to scan row", httpx.BadRequest)
// 		return
// 	}

// 	body.ID = id
// 	body.Status = status
// 	body.CreatedAt = created_at
// 	body.UpdatedAt = &updated_at

// 	lh.Logger.Info("listing created successfully", "request_id", requestId, "listing_id", id)

// 	httpx.Write(w, http.StatusCreated, body)
// }
