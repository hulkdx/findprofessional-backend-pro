package booking

import (
	"context"

	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetPriceAndCurrency(ctx context.Context, proId string) (int64, string, error)
	InsertBooking(ctx context.Context, userId int64, proId string, req *booking_model.CreateBookingRequest) (int64, error)
}

type repositoryImpl struct {
	db           *pgxpool.Pool
	timeProvider utils.TimeProvider
}

func NewRepository(db *pgxpool.Pool, timeProvider utils.TimeProvider) Repository {
	return &repositoryImpl{
		db:           db,
		timeProvider: timeProvider,
	}
}

func (r *repositoryImpl) GetPriceAndCurrency(ctx context.Context, proId string) (int64, string, error) {
	var priceNumber int64
	var priceCurrency string
	row := r.db.QueryRow(ctx, `SELECT price_number, price_currency FROM professionals WHERE id = $1`, proId)
	err := row.Scan(&priceNumber, &priceCurrency)
	if err != nil {
		return 0, "", err
	}
	return priceNumber, priceCurrency, nil
}

func (r *repositoryImpl) InsertBooking(ctx context.Context, userId int64, proId string, req *booking_model.CreateBookingRequest) (int64, error) {
	return 0, nil
}
