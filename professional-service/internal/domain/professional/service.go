package professional

import (
	"context"
	"errors"
	"fmt"
)

const REVIEW_LIMIT = 3

var (
	filterQuery = fmt.Sprintf(`
	p.id,
	p.email,
	p.first_name,
	p.last_name,
	p.coach_type,
	p.price_number,
	p.price_currency,
	p.profile_image_url,
	p.description,
	AVG(r.rate)::numeric(10,2) AS rating,
	COUNT(r),
	jsonb_agg(a) FILTER (WHERE a.id IS NOT NULL),
	jsonb_agg(json_build_object(
		'id', r.id,
		'rate', r.rate,
		'contentText', r.content_text,
		'createdAt', r.created_at,
		'updatedAt', r.updated_at,
		'user', json_build_object(
			'id', u.id,
			'email', u.email,
			'firstName', u.first_name,
			'lastName', u.last_name,
			'profileImage', u.profile_image
		)
		)) FILTER (WHERE r.id IS NOT NULL AND r.row_num <= %d)
		`,
		REVIEW_LIMIT,
	)
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
			&pro.ReviewSize,
			&pro.Availability,
			&pro.Review,
		}
	}
)

var ErrNotFound = errors.New("not found")

type Service interface {
	FindAll(context.Context) ([]Professional, error)
	FindById(ctx context.Context, id string) (Professional, error)
	Create(context.Context, CreateRequest) error
	Update(ctx context.Context, id string, p UpdateRequest) error
	FindAllReview(ctx context.Context, professionalId int64, page int, pageSize int) (Reviews, error)
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

func (s *serviceImpl) Create(ctx context.Context, r CreateRequest) error {
	return s.repository.Create(ctx, r)
}

func (s *serviceImpl) Update(ctx context.Context, id string, p UpdateRequest) error {
	return s.repository.Update(ctx, id, p)
}

func (s *serviceImpl) FindAllReview(ctx context.Context, professionalId int64, page int, pageSize int) (Reviews, error) {
	return s.repository.FindAllReview(ctx, professionalId, page, pageSize)
}
