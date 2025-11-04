package booking

import (
	"context"
	"time"

	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/payment"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

type Service struct {
	repository     Repository
	paymentService payment.PaymentService
	timeProvider   utils.TimeProvider
}

func NewService(repository Repository, paymentService payment.PaymentService, timeProvider utils.TimeProvider) *Service {
	return &Service{
		repository:     repository,
		paymentService: paymentService,
		timeProvider:   timeProvider,
	}
}

type CreateParams struct {
	Availabilities []booking_model.Availability
	IdempotencyKey string
	AmountInCents  int64
	Currency       string
	UserId         int64
	ProId          int64
	Auth           string
}

func (s *Service) Create(ctx context.Context, params *CreateParams) (*booking_model.CreateBookingResponse, error) {
	return s.repository.WithTx(ctx, func() (*booking_model.CreateBookingResponse, error) {
		return s.create(ctx, params)
	})
}

func (s *Service) create(ctx context.Context, params *CreateParams) (*booking_model.CreateBookingResponse, error) {
	expiry := s.timeProvider.Now().UTC().Add(60 * time.Second)

	_, err := s.repository.InsertBookingHolds(ctx, params.UserId, params.IdempotencyKey, expiry)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
