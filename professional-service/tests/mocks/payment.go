package mocks

import (
	"context"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/payment"
)

type FakePaymentService struct{}

func (s *FakePaymentService) CreatePaymentIntent(ctx context.Context, userId int64, amountInCents int64, currency string) (*payment.PaymentIntentResponse, error) {
	return &payment.PaymentIntentResponse{}, nil
}
