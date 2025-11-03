package booking

import (
	"context"
	"time"

	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
)

type WithTxFunc = func() (*booking_model.CreateBookingResponse, error)

type Repository interface {
	WithTx(ctx context.Context, fn WithTxFunc) (*booking_model.CreateBookingResponse, error)
	InsertBookingHolds(ctx context.Context, UserId int64, IdempotencyKey string, expiry time.Time) (*int64, error)
}
