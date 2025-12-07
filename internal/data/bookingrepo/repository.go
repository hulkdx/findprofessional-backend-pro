package bookingrepo

import (
	"context"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryImpl struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) booking.Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) GetStatus(ctx context.Context, bookingId string) (string, error) {
	query := `
		SELECT status
		FROM bookings
		WHERE id = $1
	`
	var status string
	err := r.db.QueryRow(ctx, query, bookingId).Scan(&status)
	return status, err
}
