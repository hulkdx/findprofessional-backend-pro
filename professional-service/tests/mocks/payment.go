package mocks

import (
	"context"

	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
)

type FakePaymentService struct{}

func (s *FakePaymentService) CreatePaymentIntent(ctx context.Context, userId int64, amountInCents int64, currency string, auth string) (*booking_model.PaymentIntentResponse, error) {
	return &booking_model.PaymentIntentResponse{}, nil
}
