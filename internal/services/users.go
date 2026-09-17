package services

import (
	"context"
	"fmt"

	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"

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
	hashedPassword, err := hashPassword(password)
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

func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Tune these parameters for your server.
	time := uint32(1)
	memory := uint32(64 * 1024) // 64 MiB
	threads := uint8(4)
	keyLen := uint32(32)

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		threads,
		keyLen,
	)

	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		memory,
		time,
		threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}
