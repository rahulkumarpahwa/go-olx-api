package handlers

import (
	"time"

	"github.com/google/uuid"
	"github.com/rahulkumarpahwa/go-olx-api/internal/types"
)

type CreateListingRequest struct {
	ID          uuid.UUID           `json:"-"`
	Title       string              `json:"title"`
	Description *string             `json:"description,omitempty"`
	Price       int64               `json:"price"`
	City        string              `json:"city"`
	Status      types.ListingStatus `json:"-"`
	UserID      uuid.UUID           `json:"user_id"`
	CategoryID  uuid.UUID           `json:"category_id"`
	CreatedAt   time.Time           `json:"-"`
	UpdatedAt   *time.Time          `json:"-"`
}