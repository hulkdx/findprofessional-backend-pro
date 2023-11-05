package professional

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type FilterItems func(pro *Professional) []any

type Repository interface {
	FindAll(ctx context.Context, filterQuery string, filterItems FilterItems) ([]Professional, error)
	FindById(ctx context.Context, id string, filterQuery string, filterItems FilterItems) (Professional, error)
	Update(ctx context.Context, id string, p UpdateRequest) error
}

type repositoryImpl struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repositoryImpl{db}
}

func (r *repositoryImpl) FindAll(ctx context.Context, filterQuery string, filterItems FilterItems) ([]Professional, error) {
	query := fmt.Sprintf(`
	SELECT %s FROM professionals p
	LEFT JOIN professional_rating r
		ON p.id=r.professional_id
	LEFT JOIN professional_availability a
		ON p.id=a.professional_id
	GROUP BY p.id
	`, filterQuery)
	return r.find(ctx, filterItems, query)
}

func (r *repositoryImpl) FindById(ctx context.Context, id string, filterQuery string, filterItems FilterItems) (Professional, error) {
	query := fmt.Sprintf(`
	SELECT %s FROM professionals p
	LEFT JOIN professional_rating r
		ON p.id=r.professional_id
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
