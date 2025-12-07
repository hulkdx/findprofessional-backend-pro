package professionalrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
)

const ReviewLimit = 3

type FilterItems func(pro *professional.Professional) []any

func (r *RepositoryImpl) FindAll(ctx context.Context) ([]professional.Professional, error) {
	query := r.buildFindQuery("")
	// log.Fatal(query)
	return r.find(ctx, query)
}

func (r *RepositoryImpl) FindById(ctx context.Context, id string) (professional.Professional, error) {
	query := r.buildFindQuery("AND p.id = $1")
	return r.findOne(ctx, query, id)
}

func (r *RepositoryImpl) buildFindQuery(whereClause string) string {
	return fmt.Sprintf(`
	WITH professional_review_cte AS
	(
		SELECT
			*, 
			ROW_NUMBER() OVER (PARTITION BY professional_id ORDER BY created_at) AS row_num
		FROM professional_review
	),
	review_agg AS (
		SELECT
			r.professional_id,
			AVG(r.rate)::numeric(10,2) AS rating,
			COUNT(*)                    AS review_count,
			jsonb_agg(
				json_build_object(
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
				)
				ORDER BY r.created_at
			) FILTER (WHERE r.row_num <= %d) AS reviews_json
		FROM professional_review_cte r
		
		LEFT JOIN users u
			ON u.id = r.user_id
		
		GROUP BY r.professional_id
	),
	availability_agg AS (
		SELECT
			a.professional_id,
			jsonb_agg(
				json_build_object(
					'id', a.id,
					'from', lower(a.availability),
					'to', upper(a.availability),
					'createdAt', a.created_at,
					'updatedAt', a.updated_at
				)
			) FILTER (WHERE a.id IS NOT NULL AND bi.id IS NULL) AS availability_json
		FROM professional_availability a

		LEFT JOIN booking_items bi
			ON bi.availability_id = a.id
		
		WHERE lower(a.availability) > '%s'
		
		GROUP BY a.professional_id
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
    	ra.rating AS rating,
    	COALESCE(ra.review_count, 0),
		aa.availability_json,
		ra.reviews_json
	FROM professionals p
	
	-- review
	LEFT JOIN review_agg ra
		ON ra.professional_id = p.id
	
	-- availability
	LEFT JOIN availability_agg aa
		ON p.id=aa.professional_id

	WHERE p.price_currency IS NOT NULL AND
		  p.price_number   IS NOT NULL AND
          p.pending        = false
		  %s
	`,
		ReviewLimit,
		r.timeProvider.Now().UTC().Format("2006-01-02 15:04:05"),
		whereClause,
	)
}

func scanFindFunc(pro *professional.Professional) []any {
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

func (r *RepositoryImpl) find(ctx context.Context, query string, args ...any) ([]professional.Professional, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	professionals := []professional.Professional{}
	for rows.Next() {
		pro := professional.Professional{
			Availability: []professional.Availability{},
			Review:       []professional.Review{},
		}
		err := rows.Scan(scanFindFunc(&pro)...)
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

func (r *RepositoryImpl) findOne(ctx context.Context, query string, queryArgs ...any) (professional.Professional, error) {
	professionals, err := r.find(ctx, query, queryArgs...)
	if err != nil {
		return professional.Professional{}, err
	}
	if len(professionals) == 0 {
		return professional.Professional{}, sql.ErrNoRows
	}
	return professionals[0], nil
}
