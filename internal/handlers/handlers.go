package handlers

import "github.com/rahulkumarpahwa/go-olx-api/internal/services"

type Handlers struct {
	HealthServices  *services.HealthServices
	UserServices    *services.UserServices
	ListingServices *services.ListingServices
}

func NewUserHandlers(userServices *services.UserServices) *Handlers {
	return &Handlers{
		UserServices: userServices,
	}
}

func NewListingHanlders(listingServices *services.ListingServices) *Handlers {
	return &Handlers{
		ListingServices: listingServices,
	}
}

func NewHealthHandlers(healthServices *services.HealthServices) *Handlers {
	return &Handlers{
		HealthServices: healthServices,
	}
}
