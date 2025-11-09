package bookingrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/jackc/pgx/v5"
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

func (r *repositoryImpl) WithTx(ctx context.Context, fn booking.WithTxFunc) (*int64, error) {
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

	query := `INSERT INTO booking_holds (user_id, idempotency_key, created_at, expires_at) VALUES ($1, $2, $3, $4) 
       ON CONFLICT (user_id, idempotency_key) DO NOTHING
       RETURNING id;`
	err := r.tx.QueryRow(ctx, query, UserId, IdempotencyKey, createdAt, expiry).Scan(&holdId)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, utils.ErrIdempotencyKeyIsUsed
	}
	if err != nil {
		return nil, err
	}

	return &holdId, nil
}

func (r *repositoryImpl) GetBookingHold(ctx context.Context, userId int64, idempotencyKey string) (*bookingmodel.BookingHold, error) {
	query := `SELECT (id, user_id, idempotency_key, created_at, expires_at) FROM booking_holds WHERE user_id = $1 AND idempotency_key = $2;`
	row := r.tx.QueryRow(ctx, query, userId, idempotencyKey)

	var bookingHold bookingmodel.BookingHold
	if err := row.Scan(&bookingHold); err != nil {
		return nil, err
	}
	return &bookingHold, nil
}

func (r *repositoryImpl) InsertBookingHoldItems(
	ctx context.Context,
	holdId int64,
	availabilities []bookingmodel.Availability,
	expiry time.Time,
	professionalId int64,
) error {
	now := r.timeProvider.Now()
	rows := make([][]any, len(availabilities))
	ids := make([]int64, len(availabilities))
	for i, a := range availabilities {
		rows[i] = []any{
			holdId,
			a.Id,
			now,
			expiry,
		}
		ids[i] = a.Id
	}

	err := r.ensureAvailabilitiesBelongToProfessional(ctx, ids, professionalId)
	if err != nil {
		return err
	}

	count, err := r.tx.CopyFrom(
		ctx,
		pgx.Identifier{"booking_hold_items"},
		[]string{"hold_id", "availability_id", "created_at", "expires_at"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}
	if count != int64(len(availabilities)) {
		return sql.ErrNoRows
	}

	return nil
}

func (r *repositoryImpl) EnsureAvailabilitiesBelongToProfessional(
	ctx context.Context,
	availabilities []bookingmodel.Availability,
	professionalId int64,
) error {
	ids := make([]int64, len(availabilities))
	for i, a := range availabilities {
		ids[i] = a.Id
	}
	return r.ensureAvailabilitiesBelongToProfessional(ctx, ids, professionalId)
}

func (r *repositoryImpl) ensureAvailabilitiesBelongToProfessional(
	ctx context.Context,
	ids []int64,
	professionalId int64,
) error {
	query := `SELECT EXISTS (
		SELECT id FROM professional_availability WHERE id = ANY($1::bigint[]) AND professional_id <> $2)`
	var hasMismatch bool
	err := r.tx.QueryRow(ctx, query, ids, professionalId).Scan(&hasMismatch)
	if err != nil {
		return err
	}
	if hasMismatch {
		return utils.ErrAvailabilityOwnershipMismatch
	}
	return nil
}
