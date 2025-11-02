package booking

import (
	"context"

	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/payment"
)

type Service struct {
	repository     Repository
	paymentService payment.PaymentService
}

func NewService(repository Repository, paymentService payment.PaymentService) *Service {
	return &Service{
		repository:     repository,
		paymentService: paymentService,
	}
}

type CreateParams struct {
	Availabilities []booking_model.Availability
	IdempotencyKey string
	AmountInCents  int64
	Currency       string
	UserId         int64
	ProId          string
	Auth           string
}

func (s *Service) Create(ctx context.Context, params *CreateParams) (*booking_model.CreateBookingResponse, error) {
	return nil, nil
}
