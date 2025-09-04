package booking

import "context"

type BookingService struct {
	repository Repository
}

func NewService(repository Repository) *BookingService {
	return &BookingService{
		repository: repository,
	}
}

func (s *BookingService) Create(ctx context.Context, professionalID string, req CreateBookingRequest) (*CreateBookingResponse, error) {
	return &CreateBookingResponse{}, nil
}
