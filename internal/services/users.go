package services

import (
	"context"
	"fmt"

	"github.com/rahulkumarpahwa/go-olx-api/internal/auth"
	"github.com/rahulkumarpahwa/go-olx-api/internal/dto/users"
	"github.com/rahulkumarpahwa/go-olx-api/internal/repositories"
)

// todo : add logger here as well
type UserServices struct {
	Storage repositories.UserStorage
}

func NewUserService(storage repositories.UserStorage) *UserServices {
	return &UserServices{
		Storage: storage,
	}
}

func (u *UserServices) Signup(ctx context.Context, body users.CreateUser, requestId string) (string, error) {

	// getting user by email
	data, err := u.Storage.GetUserByEmail(ctx, requestId, body.Email)
	if err != nil {
		return "", err
	}

	// check if the user exists
	if data.Email == body.Email {
		return "", fmt.Errorf("user exists already")
	}

	// password hashing
	password := body.Password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return "", err
	}
	body.Password = hashedPassword

	id, err := u.Storage.CreateUser(ctx, requestId, body)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
