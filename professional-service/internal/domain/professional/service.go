package professional

import (
	"context"
	"errors"
)

var (
	filterQuery = `
	p.id,
	p.email,
	p.first_name,
	p.last_name,
	p.coach_type,
	p.price_number,
	p.price_currency,
	p.profile_image_url,
	p.description,
	AVG(rate)::numeric(10,2) AS rating,
	jsonb_agg(a) FILTER (WHERE a.id IS NOT NULL),
	jsonb_agg(json_build_object(
		'id', r.id,
		'rate', r.rate,
		'contentText', r.content_text,
		'user', json_build_object(
			'id', u.id,
			'email', u.email,
			'firstName', u.first_name,
			'lastName', u.last_name,
			'profileImage', u.profile_image
		)
		)) FILTER (WHERE r.id IS NOT NULL)
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
			&pro.Review,
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
