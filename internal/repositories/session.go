package repositories

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/rahulkumarpahwa/go-olx-api/internal/config"
	"github.com/rahulkumarpahwa/go-olx-api/internal/dto/session"
)

type SessionStorage interface {
	CreateSession(ctx context.Context, userID uuid.UUID, tokenHash string, expiry time.Time) (*session.CreateSession, error)
	UpdateSession()
}

type SessionRepositories struct {
	Config *config.Config
	DB     *sql.DB
	Logger *slog.Logger
}

func NewSessionRepository(cfg *config.Config, db *sql.DB, logger *slog.Logger) *SessionRepositories {
	return &SessionRepositories{
		Config: cfg,
		DB:     db,
		Logger: logger,
	}
}

func (sr *SessionRepositories) CreateSession(ctx context.Context, userID uuid.UUID, tokenHash string, expiry time.Time) (*session.CreateSession, error) {

	const query = `INSERT INTO refresh_sessions 
	(user_id, token_hash, expires_at) 
	VALUES ($1, $2, $3)
	RETURNING id;`

	row := sr.DB.QueryRowContext(ctx, query, userID, tokenHash, expiry)
	err := row.Err()
	if err != nil {
		sr.Logger.Error("insert session row error", "user_id", userID, "err", err)
		return nil, err
	}

	refToken := session.CreateSession{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiry,
	}

	err = row.Scan(&refToken.ID)
	if err != nil {
		sr.Logger.Error("scanned session row error", "user_id", userID, "err", err)
		return nil, err
	}

	return &refToken, nil
}

// func (sr *SessionRepositories) UpdateSession() {
// 	const query = `UPDATE refresh_sessions
// 					SET token_hash = $1,
//     				expires_at = $2,
//     				updated_at = NOW()
// 					WHERE id = $3
//   				AND revoked_at IS NULL;`

	


// }
