package handlers

import (
	"fmt"
	"strings"
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

type ValidationError struct {
	Field string
	Msg   string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf(`%s: %s`, v.Field, v.Msg)
}

func (l *CreateListingRequest) Validate() error {
	if strings.TrimSpace(l.Title) == "" {
		return &ValidationError{Field: "title", Msg: "must not be empty"}
	}

	if len(l.Title) > 200 {
		return &ValidationError{Field: "title", Msg: "must be at most 200 characters"}
	}

	if l.Price <= 0 {
		return &ValidationError{Field: "price", Msg: "must be greater than 0"}
	}

	if strings.TrimSpace(l.City) == "" {
		return &ValidationError{Field: "city", Msg: "must not be empty"}
	}

	if l.UserID == uuid.Nil {
		return &ValidationError{Field: "user_id", Msg: "must not be empty"}
	}

	if l.CategoryID == uuid.Nil {
		return &ValidationError{Field: "catergory_id", Msg: "must not be empty"}
	}

	return nil
}
