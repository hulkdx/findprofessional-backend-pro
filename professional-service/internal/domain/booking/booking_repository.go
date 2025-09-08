package booking

import (
	"context"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetPriceAndCurrency(ctx context.Context, proId string) (int64, string, error)
	InsertBooking(ctx context.Context, userId int64, proId string, req CreateBookingRequest) error
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

func (r *repositoryImpl) InsertBooking(ctx context.Context, userId int64, proId string, req CreateBookingRequest) error {
	status := BookingStatusHold

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	txDone := false
	defer func() {
		if txDone {
			tx.Commit(ctx)
		} else {
			tx.Rollback(ctx)
		}
	}()

	query := `
		INSERT INTO bookings (
			user_id,
			professional_id,
			status,
			amount_in_cents,
			currency,
			idempotency_key,
			created_at,
			updated_at,
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
	`
	row := tx.QueryRow(ctx, query,
		userId,
		proId,
		status,
		req.AmountInCents,
		req.Currency,
		req.IdempotencyKey,
		r.timeProvider.Now(),
		r.timeProvider.Now(),
	)
	var bookingId int64
	err = row.Scan(&bookingId)
	if err != nil {
		return err
	}

	rows := make([][]any, len(req.Slots))
	now := r.timeProvider.Now()
	for i, slot := range req.Slots {
		tsRange, err := utils.ConvertToTsRange(slot.Date, slot.From, slot.To)
		if err != nil {
			return err
		}
		rows[i] = []any{
			bookingId,
			proId,
			tsRange,
			status,
			now,
			now,
		}
	}

	columns := []string{
		"booking_id",
		"professional_id",
		"slot",
		"status",
		"created_at",
		"updated_at",
	}
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"booking_slots"},
		columns,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	txDone = true
	return nil
}
