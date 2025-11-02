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

func (s *Service) Create(ctx context.Context, userId int64, proId string, req *booking_model.CreateBookingRequest, auth string) (*booking_model.CreateBookingResponse, error) {
	return nil, nil
}
