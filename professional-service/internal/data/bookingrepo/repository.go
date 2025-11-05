package bookingrepo

import (
	"context"
	"errors"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repositoryImpl struct {
	db           *pgxpool.Pool
	timeProvider utils.TimeProvider
	tx           pgx.Tx
}

func NewRepository(db *pgxpool.Pool, timeProvider utils.TimeProvider) booking.Repository {
	return &repositoryImpl{
		db:           db,
		timeProvider: timeProvider,
	}
}

func (r *repositoryImpl) WithTx(ctx context.Context, fn booking.WithTxFunc) (*bookingmodel.CreateBookingResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	r.tx = tx

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

func (r *repositoryImpl) InsertBookingHolds(ctx context.Context, UserId int64, IdempotencyKey string, expiry time.Time) (*int64, error) {
	createdAt := r.timeProvider.Now()
	var holdId int64

	query := `INSERT INTO booking_holds (user_id, idempotency_key, created_at, expires_at) VALUES ($1, $2, $3, $4) RETURNING id;`
	err := r.tx.QueryRow(ctx, query, UserId, IdempotencyKey, createdAt, expiry).Scan(&holdId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation && pgErr.ConstraintName == "booking_holds_user_ik_uk" {
			return nil, utils.ErrIdempotencyKeyIsUsed
		}
		return nil, err
	}

	return &holdId, nil
}

func (r *repositoryImpl) GetBookingHold(ctx context.Context, userId int64, idempotencyKey string) (*bookingmodel.BookingHold, error) {
	query := `SELECT * FROM booking_holds WHERE user_id = $1 AND idempotency_key = $2;`
	row := r.tx.QueryRow(ctx, query, userId, idempotencyKey)

	var holdId bookingmodel.BookingHold
	if err := row.Scan(&holdId); err != nil {
		return nil, err
	}
	return &holdId, nil
}
