package booking

import (
	"context"
	"errors"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/payment"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

type Service struct {
	repository     Repository
	paymentService payment.Service
	timeProvider   utils.TimeProvider
}

func NewService(repository Repository, paymentService payment.Service, timeProvider utils.TimeProvider) *Service {
	return &Service{
		repository:     repository,
		paymentService: paymentService,
		timeProvider:   timeProvider,
	}
}

type CreateParams struct {
	Availabilities []bookingmodel.Availability
	IdempotencyKey string
	AmountInCents  int64
	Currency       string
	UserId         int64
	ProId          int64
	Auth           string
}

func (s *Service) Create(ctx context.Context, params *CreateParams) (*bookingmodel.CreateBookingResponse, error) {
	// TODO: validation
	return s.repository.WithTx(ctx, func() (*bookingmodel.CreateBookingResponse, error) {
		return s.create(ctx, params)
	})
}

func (s *Service) create(ctx context.Context, params *CreateParams) (*bookingmodel.CreateBookingResponse, error) {
	expiry := s.timeProvider.Now().UTC().Add(60 * time.Second)
	holdId, err := s.getBookingHoldId(ctx, params, expiry)
	if err != nil {
		return nil, err
	}
	payResponse, err := s.paymentService.CreatePaymentIntent(
		ctx,
		holdId,
		params.AmountInCents,
		params.Currency,
		params.IdempotencyKey,
		params.Auth,
	)
	if err != nil {
		return nil, err
	}

	return &bookingmodel.CreateBookingResponse{
		PaymentIntentResponse: *payResponse,
	}, nil
}

func (s *Service) getBookingHoldId(ctx context.Context, params *CreateParams, expiry time.Time) (int64, error) {
	holdId, err := s.repository.InsertBookingHolds(ctx, params.UserId, params.IdempotencyKey, expiry)
	if errors.Is(err, utils.ErrIdempotencyKeyIsUsed) {
		hold, err1 := s.repository.GetBookingHold(ctx, params.UserId, params.IdempotencyKey)
		if err1 != nil {
			return -1, errors.Join(err, err1)
		}
		return hold.ID, nil
	}
	if err != nil {
		return -1, err
	}
	return *holdId, nil
}
