package repositories

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/rahulkumarpahwa/go-olx-api/internal/dto"
)

type UserStorage interface {
	GetUserById(ctx context.Context, requestId int, id string) (dto.RequestUser, error)
}

type UserRepositories struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewUserRepository(db *sql.DB, logger *slog.Logger) *UserRepositories {
	return &UserRepositories{
		DB:     db,
		Logger: logger,
	}
}

func (r *UserRepositories) GetUserById(ctx context.Context, requestId int, id string) (dto.RequestUser, error) {
	const query = "SELECT id, name, email FROM user WHERE id=$1"
	row := r.DB.QueryRowContext(ctx, query, id)

	var req dto.RequestUser

	err := row.Scan(&req.ID, &req.Email, &req.Name)

	if err != nil {
		r.Logger.Error("scanned user row error", "request_id", requestId, "err", err, "user_id", id)
		return dto.RequestUser{}, err
	}

	err = row.Err()
	if err != nil {
		r.Logger.Error("scanned user row error", "request_id", requestId, "err", err)
		return dto.RequestUser{}, err
	}

	return req, nil
}
