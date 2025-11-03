package booking

import (
	"context"

	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
)

type Func = func() (*booking_model.CreateBookingResponse, error)

type Repository interface {
	WithTx(ctx context.Context, fn Func) (*booking_model.CreateBookingResponse, error)
}
