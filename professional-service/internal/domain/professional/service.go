package professional

import (
	"context"
	"errors"
)

var (
	filterQuery = `
	p.id,
	email,
	first_name,
	last_name,
	coach_type,
	price_number,
	price_currency,
	profile_image_url,
	description,
	AVG(rate)::numeric(10,2) AS rating,
	jsonb_agg(a) FILTER (WHERE a IS NOT NULL)
`
	filterItems = func(pro *Professional) []any {
		return []any{
			&pro.ID,
			&pro.Email,
			&pro.FirstName,
			&pro.LastName,
			&pro.CoachType,
			&pro.PriceNumber,
			&pro.PriceCurrency,
			&pro.ProfileImageUrl,
			&pro.Description,
			&pro.Rating,
			&pro.Availability,
		}
	}
)

var ErrNotFound = errors.New("not found")

type Service interface {
	FindAll(context.Context) ([]Professional, error)
	FindById(ctx context.Context, id string) (Professional, error)
	Update(ctx context.Context, id string, p UpdateRequest) error
}

type serviceImpl struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &serviceImpl{repository}
}

func (s *serviceImpl) FindAll(ctx context.Context) ([]Professional, error) {
	return s.repository.FindAll(ctx, filterQuery, filterItems)
}

func (s *serviceImpl) FindById(ctx context.Context, id string) (Professional, error) {
	return s.repository.FindById(ctx, id, filterQuery, filterItems)
}

func (s *serviceImpl) Update(ctx context.Context, id string, p UpdateRequest) error {
	return s.repository.Update(ctx, id, p)
}
