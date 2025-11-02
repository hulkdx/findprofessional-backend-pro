package booking

import (
	"context"

	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/payment"
)

type BookingService struct {
	repository     Repository
	paymentService payment.PaymentService
}

func NewService(repository Repository, paymentService payment.PaymentService) *BookingService {
	return &BookingService{
		repository:     repository,
		paymentService: paymentService,
	}
}

func (s *BookingService) Create(ctx context.Context, userId int64, proId string, req *booking_model.CreateBookingRequest, auth string) (*booking_model.CreateBookingResponse, error) {
	return nil, nil
}
