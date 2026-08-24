package types

import "time"

type ListingStatus string

const (
	ListingStatusActive   ListingStatus = "active"
	ListingStatusInactive ListingStatus = "inactive"
)

type Listings struct {
	Id          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Price       int64         `json:"price"`
	City        string        `json:"city"`
	Status      ListingStatus `json:"status"`
	UserId      string        `json:"user_id"`
	CategoryId  string        `json:"category_id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
