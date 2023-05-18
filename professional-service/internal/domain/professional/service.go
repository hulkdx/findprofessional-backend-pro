package professional

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type Service interface {
	FindAll(context.Context) ([]Professional, error)
}

type serviceImpl struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &serviceImpl{repository}
}

func (s *serviceImpl) FindAll(ctx context.Context) ([]Professional, error) {
	return s.repository.FindAll(ctx, "ID", "Email")
}
