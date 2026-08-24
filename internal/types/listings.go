package types

import (
	"time"
	"github.com/google/uuid"
)

type ListingStatus string

const (
	ListingStatusActive   ListingStatus = "active"
	ListingStatusInactive ListingStatus = "inactive"
)

type Listings struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	Description *string       `json:"description"`
	Price       int64         `json:"price"`
	City        string        `json:"city"`
	Status      ListingStatus `json:"status"`
	UserID      uuid.UUID     `json:"user_id"`
	CategoryID  uuid.UUID     `json:"category_id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   *time.Time    `json:"updated_at"`
}
