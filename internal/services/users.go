package services

import "github.com/rahulkumarpahwa/go-olx-api/internal/repositories"

type UserServices struct {
	Storage repositories.UserStorage
}

func NewUserService(storage repositories.UserStorage) *UserServices {
	return &UserServices{
		Storage: storage,
	}
}
