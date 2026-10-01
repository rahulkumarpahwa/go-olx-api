package session

import (
	"time"

	"github.com/google/uuid"
)

type CreateSession struct {
	ID uuid.UUID
	UserID uuid.UUID
	TokenHash string
	ExpiresAt time.Time
}