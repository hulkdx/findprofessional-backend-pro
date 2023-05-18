package professional

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Repository interface {
	FindAll(ctx context.Context, fields ...string) ([]Professional, error)
	FindById(ctx context.Context, id string, fields ...string) (Professional, error)
	Update(ctx context.Context, id string, p Professional) error
}

type repositoryImpl struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repositoryImpl{db}
}

func (r *repositoryImpl) FindAll(ctx context.Context, fields ...string) ([]Professional, error) {
	query := fmt.Sprintf("SELECT %s FROM professionals", strings.Join(fields, ", "))
	return r.find(ctx, fields, query)
}

func (r *repositoryImpl) FindById(ctx context.Context, id string, fields ...string) (Professional, error) {
	query := fmt.Sprintf("SELECT %s FROM professionals WHERE id = ?", strings.Join(fields, ", "))
	return r.findOne(ctx, fields, query, id)
}

func (r *repositoryImpl) Update(ctx context.Context, id string, p Professional) error {
	query := "UPDATE professionals SET email = ?, password = ?, updated_at = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, p.Email, p.Password, p.UpdatedAt, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) find(ctx context.Context, fields []string, query string, args ...any) ([]Professional, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	professionals := []Professional{}
	for rows.Next() {
		pro := Professional{}
		elem := reflect.ValueOf(&pro).Elem()
		scanArgs := make([]interface{}, len(fields))
		for i := range fields {
			field := elem.FieldByName(fields[i])
			scanArgs[i] = field.Addr().Interface()
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}
		professionals = append(professionals, pro)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return professionals, nil
}

func (r *repositoryImpl) findOne(ctx context.Context, fields []string, query string, queryArgs ...any) (Professional, error) {
	professionals, err := r.find(ctx, fields, query, queryArgs...)
	if err != nil {
		return Professional{}, err
	}
	if len(professionals) == 0 {
		return Professional{}, sql.ErrNoRows
	}
	return professionals[0], nil
}
