package professional

import (
	"context"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Professional, error)
	FindById(ctx context.Context, id string) (Professional, error)
	Create(ctx context.Context, request CreateRequest, pending bool) error
	Update(ctx context.Context, id string, p UpdateRequest) error
	FindAllReview(ctx context.Context, professionalId int64, page int, pageSize int) (Reviews, error)
	GetAvailability(ctx context.Context, professionalId int64) (Availabilities, error)
	UpdateAvailability(ctx context.Context, professionalId int64, availability UpdateAvailabilityRequest) error
	GetBookingStatus(ctx context.Context, bookingId int64) (StatusResponse, error)
}
