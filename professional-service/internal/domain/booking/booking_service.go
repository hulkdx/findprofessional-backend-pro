package booking

import (
	"context"
	"errors"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/payment"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/logger"
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

func (s *Service) Create(ctx context.Context, params *CreateParams) (*bookingmodel.PaymentResponse, error) {
	holdId, err := s.repository.WithTx(ctx, func() (*int64, error) {
		return s.createTx(ctx, params)
	})
	if err != nil {
		return nil, err
	}

	logger.Debug("Creating payment intent for hold ID", *holdId, "amount", params.AmountInCents)
	payResponse, err := s.paymentService.CreatePaymentIntent(
		ctx,
		*holdId,
		params.AmountInCents,
		params.Currency,
		params.IdempotencyKey,
		params.Auth,
		params.ProId,
	)
	if err != nil {
		logger.Error("Failed to create payment intent", err)
	}
	return payResponse, err
}

func (s *Service) createTx(ctx context.Context, params *CreateParams) (*int64, error) {
	expiry := s.timeProvider.Now().UTC().Add(60 * time.Second)

	holdId, err := s.repository.InsertBookingHolds(ctx, params.UserId, params.IdempotencyKey, expiry)
	if err != nil {
		if errors.Is(err, utils.ErrIdempotencyKeyIsUsed) {
			return s.getBookingHold(ctx, params)
		}
		logger.Error("Failed to insert booking hold", err)
		return nil, err
	}
	err = s.repository.InsertBookingHoldItems(ctx, *holdId, params.Availabilities, expiry, params.ProId)
	if err != nil {
		logger.Error("Failed to insert booking hold items", err)
		return nil, err
	}
	return holdId, nil
}

func (s *Service) getBookingHold(ctx context.Context, params *CreateParams) (*int64, error) {
	hold, err := s.repository.GetBookingHold(ctx, params.UserId, params.IdempotencyKey)
	if err != nil {
		logger.Error("Failed to get booking hold", err)
		return nil, err
	}
	err = s.repository.EnsureAvailabilitiesBelongToProfessional(ctx, params.Availabilities, params.ProId)
	if err != nil {
		logger.Error("availabilities is not belonging to the professional id: ", err)
		return nil, err
	}
	return &hold.ID, nil
}
