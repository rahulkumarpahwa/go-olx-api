package repositories

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/rahulkumarpahwa/go-olx-api/internal/config"
	"github.com/rahulkumarpahwa/go-olx-api/internal/dto/users"
)

type UserStorage interface {
	GetUserById(ctx context.Context, requestId string, id string) (users.RequestUser, error)
	GetUserByEmail(ctx context.Context, requestId string, email string) (users.RequestUser, error)
	CreateUser(ctx context.Context, requestId string, body users.CreateUser) (uuid.UUID, error)
}

type UserRepositories struct {
	Config *config.Config
	DB     *sql.DB
	Logger *slog.Logger
}

func NewUserRepository(cfg *config.Config, db *sql.DB, logger *slog.Logger) *UserRepositories {
	return &UserRepositories{
		Config: cfg,
		DB:     db,
		Logger: logger,
	}
}

func (r *UserRepositories) GetUserById(ctx context.Context, requestId string, id string) (users.RequestUser, error) {
	const query = "SELECT id, name, email FROM user WHERE id=$1"
	row := r.DB.QueryRowContext(ctx, query, id)

	var req users.RequestUser

	err := row.Scan(&req.ID, &req.Email, &req.Name)

	if err != nil {
		r.Logger.Error("scanned user row error", "request_id", requestId, "err", err, "user_id", id)
		return users.RequestUser{}, err
	}

	err = row.Err()
	if err != nil {
		r.Logger.Error("scanned user row error", "request_id", requestId, "err", err)
		return users.RequestUser{}, err
	}

	return req, nil
}

func (r *UserRepositories) GetUserByEmail(ctx context.Context, requestId string, email string) (users.RequestUser, error) {
	const query = "SELECT id, email, name FROM user WHERE email=$1"
	row := r.DB.QueryRowContext(ctx, query, email)

	var req users.RequestUser

	err := row.Scan(&req.ID, &req.Email, &req.Name)

	if err != nil {
		r.Logger.Error("scanned user row error", "request_id", requestId, "err", err, "email", email)
		return users.RequestUser{}, err
	}

	err = row.Err()
	if err != nil {
		r.Logger.Error("scanned user row error", "request_id", requestId, "err", err)
		return users.RequestUser{}, err
	}

	return req, nil
}

func (r *UserRepositories) CreateUser(ctx context.Context, requestId string, body users.CreateUser) (uuid.UUID, error) {
	const query = "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id;"
	row := r.DB.QueryRowContext(ctx, query, body.Name, body.Email, body.Password)

	var req users.RequestUser

	err := row.Scan(&req.ID)

	if err != nil {
		r.Logger.Error("scanned user row error", "request_id", requestId, "err", err, "id", req.ID)
		return uuid.Nil, err
	}

	err = row.Err()
	if err != nil {
		r.Logger.Error("scanned user row error", "request_id", requestId, "err", err)
		return uuid.Nil, err
	}

	return req.ID, nil
}
