package services

import (
	"github.com/rahulkumarpahwa/go-olx-api/internal/repositories"
)

type ListingServices struct {
	Storage repositories.ListingStorage
}

func NewListingService(storage repositories.ListingStorage) *ListingServices {
	return &ListingServices{
		Storage: storage,
	}
}
