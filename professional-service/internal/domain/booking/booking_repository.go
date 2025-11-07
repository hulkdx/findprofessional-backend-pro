package booking

import (
	"context"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
)

type WithTxFunc = func() (*bookingmodel.CreateBookingResponse, error)

type Repository interface {
	WithTx(ctx context.Context, fn WithTxFunc) (*bookingmodel.CreateBookingResponse, error)
	InsertBookingHolds(ctx context.Context, UserId int64, IdempotencyKey string, expiry time.Time) (*int64, error)
	GetBookingHold(ctx context.Context, userId int64, idempotencyKey string) (*bookingmodel.BookingHold, error)
	InsertBookingHoldItems(ctx context.Context, holdId int64, availabilities []bookingmodel.Availability, expiry time.Time) error
}
