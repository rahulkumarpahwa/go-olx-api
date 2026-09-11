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
	Field   string
	Msg string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf(`%s:%s`, v.Field, v.Msg)
}

func (l *CreateListingRequest) Validate() error {
	if strings.TrimSpace(l.Title) == "" {
		return &ValidationError{Field: "title", Msg: "title field is required"}
	}

	return nil
}
