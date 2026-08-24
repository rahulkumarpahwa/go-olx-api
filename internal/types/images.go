package types

import "github.com/google/uuid"

type Images struct {
	ID        uuid.UUID `json:"id"`
	ListingID uuid.UUID `json:"listing_id"`
}
