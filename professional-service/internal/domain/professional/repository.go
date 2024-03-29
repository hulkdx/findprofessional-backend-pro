package professional

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

type FilterItems func(pro *Professional) []any

type Repository interface {
	FindAll(ctx context.Context, filterQuery string, filterItems FilterItems) ([]Professional, error)
	FindById(ctx context.Context, id string, filterQuery string, filterItems FilterItems) (Professional, error)
	Update(ctx context.Context, id string, p UpdateRequest) error
	FindAllReview(ctx context.Context, professionalId int64, page int, pageSize int) (Reviews, error)
}

type repositoryImpl struct {
	db           *sql.DB
	timeProvider utils.TimeProvider
}

func NewRepository(db *sql.DB, timeProvider utils.TimeProvider) Repository {
	return &repositoryImpl{
		db,
		timeProvider,
	}
}

func (r *repositoryImpl) FindAll(ctx context.Context, filterQuery string, filterItems FilterItems) ([]Professional, error) {
	query := fmt.Sprintf(`
	WITH professional_review_cte AS
	(
    SELECT
			*, 
			ROW_NUMBER() OVER (PARTITION BY professional_id ORDER BY created_at) AS row_num
    FROM professional_review
	)

	SELECT %s FROM professionals p
	LEFT JOIN professional_review_cte r
		ON p.id=r.professional_id
	LEFT JOIN users u
		ON r.user_id=u.id
	LEFT JOIN professional_availability a
		ON p.id=a.professional_id
			AND a.date > '%s'
	WHERE p.price_currency IS NOT NULL AND
				p.price_number   IS NOT NULL
	GROUP BY p.id
	`,
		filterQuery,
		r.timeProvider.Now().Format("2006-01-01"),
	)
	return r.find(ctx, filterItems, query)
}

func (r *repositoryImpl) FindById(ctx context.Context, id string, filterQuery string, filterItems FilterItems) (Professional, error) {
	query := fmt.Sprintf(`
	WITH professional_review_cte AS
	(
    SELECT
			*, 
			ROW_NUMBER() OVER (PARTITION BY professional_id ORDER BY created_at) AS row_num
    FROM professional_review
	)

	SELECT %s FROM professionals p
	LEFT JOIN professional_review_cte r
		ON p.id=r.professional_id
	LEFT JOIN users u
		ON r.user_id=u.id
	LEFT JOIN professional_availability a
		ON p.id=a.professional_id
	WHERE p.id=$1
	GROUP BY p.id
	`, filterQuery)
	return r.findOne(ctx, filterItems, query, id)
}

func (r *repositoryImpl) Update(ctx context.Context, id string, p UpdateRequest) error {
	query := "UPDATE professionals SET email = $1, updated_at = $2 WHERE id = $3"
	result, err := r.db.ExecContext(ctx, query, p.Email, time.Now(), id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repositoryImpl) FindAllReview(ctx context.Context, professionalID int64, page int, pageSize int) (Reviews, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT 
			pr.id,
			pr.rate,
			pr.content_text,
			pr.created_at,
			pr.updated_at,
			u.id,
			u.email,
			u.first_name,
			u.last_name,
			u.profile_image
		FROM professional_review pr
		LEFT JOIN users u ON pr.user_id = u.id
		WHERE pr.professional_id = $1
		ORDER BY pr.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, professionalID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make(Reviews, 0)

	for rows.Next() {
		var review Review

		err := rows.Scan(
			&review.ID,
			&review.Rate,
			&review.ContentText,
			&review.CreatedAt,
			&review.UpdatedAt,
			&review.User.ID,
			&review.User.Email,
			&review.User.FirstName,
			&review.User.LastName,
			&review.User.ProfileImage,
		)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
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
