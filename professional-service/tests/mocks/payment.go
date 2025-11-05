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
	IdempotencyKey string,
	auth string,
) (*bookingmodel.PaymentIntentResponse, error) {
	return &bookingmodel.PaymentIntentResponse{}, nil
}
