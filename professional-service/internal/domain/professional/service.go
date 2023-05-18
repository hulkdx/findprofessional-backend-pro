package professional

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type Service interface {
	FindAll(context.Context) ([]Professional, error)
	FindById(ctx context.Context, id string, fields ...string) (Professional, error)
	Create(ctx context.Context, p Professional) error
	Delete(ctx context.Context, id string) error
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

func (s *serviceImpl) FindById(ctx context.Context, id string, fields ...string) (Professional, error) {
	return s.repository.FindById(ctx, id, fields...)
}

func (s *serviceImpl) Create(ctx context.Context, p Professional) error {
	return s.repository.Create(ctx, p)
}

func (s *serviceImpl) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *serviceImpl) Update(ctx context.Context, id string, p Professional) error {
	return s.repository.Update(ctx, id, p)
}
