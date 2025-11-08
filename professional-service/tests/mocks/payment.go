package mocks

import (
	"context"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
)

type FakePaymentService struct{}

func (s *FakePaymentService) CreatePaymentIntent(
	ctx context.Context,
	holdId int64,
	AmountInCents int64,
	Currency string,
	idempotencyKey string,
	auth string,
	professionalId int64,
) (*bookingmodel.PaymentIntentResponse, error) {
	return &bookingmodel.PaymentIntentResponse{}, nil
}
