package mocks

import (
	"context"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
)

type FakePaymentService struct {
	CreatePaymentIntentSuccess *bookingmodel.PaymentResponse
	CreatePaymentIntentError   error
}

func (s *FakePaymentService) CreatePaymentIntent(
	ctx context.Context,
	holdId int64,
	AmountInCents int64,
	Currency string,
	idempotencyKey string,
	auth string,
	professionalId int64,
) (*bookingmodel.PaymentResponse, error) {
	return s.CreatePaymentIntentSuccess, s.CreatePaymentIntentError
}
