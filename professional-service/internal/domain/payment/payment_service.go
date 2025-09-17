package payment

import "context"

type PaymentService interface {
	CreatePaymentIntent(ctx context.Context, userId int64, amountInCents int64, currency string) error
}

func NewService() PaymentService {
	return &paymentServiceImpl{}
}

type paymentServiceImpl struct {
}

func (s *paymentServiceImpl) CreatePaymentIntent(ctx context.Context, userId int64, amountInCents int64, currency string) error {
	return nil
}
