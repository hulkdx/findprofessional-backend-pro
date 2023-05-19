package professional

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type Service interface {
	FindAll(context.Context) ([]Professional, error)
	FindById(ctx context.Context, id string) (Professional, error)
	Update(ctx context.Context, id string, p Professional) error
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

func (s *serviceImpl) FindById(ctx context.Context, id string) (Professional, error) {
	return s.repository.FindById(ctx, id, "ID", "Email")
}

func (s *serviceImpl) Update(ctx context.Context, id string, p Professional) error {
	return s.repository.Update(ctx, id, p)
}
