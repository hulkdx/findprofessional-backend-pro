package professional

import (
	"context"
	"database/sql"
	"fmt"
)

const REVIEW_LIMIT = 3

type FilterItems func(pro *Professional) []any

func (r *repositoryImpl) FindAll(ctx context.Context) ([]Professional, error) {
	query := fmt.Sprintf(`
	WITH professional_review_cte AS
	(
		SELECT
			*, 
			ROW_NUMBER() OVER (PARTITION BY professional_id ORDER BY created_at) AS row_num
		FROM professional_review
	)

	SELECT
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
	FROM professionals p
	
	-- review
	LEFT JOIN professional_review_cte r
		ON p.id=r.professional_id
	LEFT JOIN users u
		ON r.user_id=u.id
	
	-- availability
	LEFT JOIN professional_availability a
		ON p.id=a.professional_id
		AND a.date > '%s'

	WHERE p.price_currency IS NOT NULL AND
				p.price_number   IS NOT NULL
	
	GROUP BY p.id
	`,
		REVIEW_LIMIT,
		r.timeProvider.Now().Format("2006-01-01"),
	)
	filterItems := func(pro *Professional) []any {
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
	return r.find(ctx, filterItems, query)
}

func (r *repositoryImpl) FindById(ctx context.Context, id string) (Professional, error) {
	query := fmt.Sprintf(`
	WITH professional_review_cte AS
	(
    SELECT
			*, 
			ROW_NUMBER() OVER (PARTITION BY professional_id ORDER BY created_at) AS row_num
    FROM professional_review
	)

	SELECT
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
	FROM professionals p
	LEFT JOIN professional_review_cte r
		ON p.id=r.professional_id
	LEFT JOIN users u
		ON r.user_id=u.id
	LEFT JOIN professional_availability a
		ON p.id=a.professional_id
	WHERE p.id=$1
	GROUP BY p.id
	`, REVIEW_LIMIT)

	filterItems := func(pro *Professional) []any {
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

	return r.findOne(ctx, filterItems, query, id)
}

func (r *repositoryImpl) find(ctx context.Context, filterItems FilterItems, query string, args ...any) ([]Professional, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	professionals := []Professional{}
	for rows.Next() {
		pro := Professional{
			Availability: []Availability{},
			Review:       []Review{},
		}
		err := rows.Scan(filterItems(&pro)...)
		if err != nil {
			return nil, err
		}
		professionals = append(professionals, pro)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return professionals, nil
}

func (r *repositoryImpl) findOne(ctx context.Context, filterItems FilterItems, query string, queryArgs ...any) (Professional, error) {
	professionals, err := r.find(ctx, filterItems, query, queryArgs...)
	if err != nil {
		return Professional{}, err
	}
	if len(professionals) == 0 {
		return Professional{}, sql.ErrNoRows
	}
	return professionals[0], nil
}
