package booking

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetPriceAndCurrency(ctx context.Context, proId string) (int64, string, error)
}

type repositoryImpl struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repositoryImpl{
		db: db,
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
