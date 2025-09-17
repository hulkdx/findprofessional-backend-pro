package mocks

import (
	"context"
)

type FakePaymentService struct{}

func (s *FakePaymentService) CreatePaymentIntent(ctx context.Context, userId int64, amountInCents int64, currency string) error {
	return nil
}
