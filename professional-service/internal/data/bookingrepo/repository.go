package bookingrepo

import (
	"context"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking"
	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repositoryImpl struct {
	db           *pgxpool.Pool
	timeProvider utils.TimeProvider
	tx           *pgx.Tx
}

func NewRepository(db *pgxpool.Pool, timeProvider utils.TimeProvider) booking.Repository {
	return &repositoryImpl{
		db:           db,
		timeProvider: timeProvider,
	}
}

func (r *repositoryImpl) WithTx(ctx context.Context, fn booking.Func) (*booking_model.CreateBookingResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	r.tx = &tx

	txDone := false
	defer func() {
		r.tx = nil
		if txDone {
			_ = tx.Commit(ctx)
		} else {
			_ = tx.Rollback(ctx)
		}
	}()

	res, err := fn()
	if err != nil {
		return nil, err
	}

	txDone = true
	return res, nil
}
