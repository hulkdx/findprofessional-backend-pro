package professional

import (
	"context"
	"database/sql"
	"time"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Professional, error)
	FindById(ctx context.Context, id string) (Professional, error)
	Update(ctx context.Context, id string, p UpdateRequest) error
}

type repositoryImpl struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repositoryImpl{db}
}

func (r *repositoryImpl) FindAll(ctx context.Context) ([]Professional, error) {
	query := "SELECT id, email FROM professionals"
	return r.find(ctx, query)
}

func (r *repositoryImpl) FindById(ctx context.Context, id string) (Professional, error) {
	query := "SELECT id, email FROM professionals WHERE id = $1"
	return r.findOne(ctx, query, id)
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

func (r *repositoryImpl) find(ctx context.Context, query string, args ...any) ([]Professional, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	professionals := []Professional{}
	for rows.Next() {
		pro := Professional{}
		if err := rows.Scan(&pro.ID, &pro.Email); err != nil {
			return nil, err
		}
		professionals = append(professionals, pro)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return professionals, nil
}

func (r *repositoryImpl) findOne(ctx context.Context, query string, queryArgs ...any) (Professional, error) {
	professionals, err := r.find(ctx, query, queryArgs...)
	if err != nil {
		return Professional{}, err
	}
	if len(professionals) == 0 {
		return Professional{}, sql.ErrNoRows
	}
	return professionals[0], nil
}
