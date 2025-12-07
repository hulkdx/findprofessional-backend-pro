package booking

import (
	"context"
)

type Service interface {
	GetStatus(ctx context.Context, bookingId string) (StatusResponse, error)
}

type serviceImpl struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &serviceImpl{
		repository: repository,
	}
}

func (s *serviceImpl) GetStatus(ctx context.Context, bookingId string) (StatusResponse, error) {
	status, err := s.repository.GetStatus(ctx, bookingId)
	if err != nil {
		return StatusResponse{}, err
	}
	return StatusResponse{Status: status}, nil
}
