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
	Availabilities   []bookingmodel.Availability
	IdempotencyKey   string
	AmountInCents    int64
	Currency         string
	UserId           int64
	ProId            int64
	Auth             string
	StripeApiVersion string
}

func (s *Service) Create(ctx context.Context, params *CreateParams) (*bookingmodel.CreateBookingResponse, error) {
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
		return nil, err
	}

	logger.Debug("Booking creation completed successfully for user", params.UserId)
	return &bookingmodel.CreateBookingResponse{
		PaymentIntentResponse: *payResponse,
	}, nil
}

func (s *Service) createTx(ctx context.Context, params *CreateParams) (*int64, error) {
	expiry := s.timeProvider.Now().UTC().Add(60 * time.Second)

	logger.Debug("Starting booking creation for user", params.UserId, "with idempotency key", params.IdempotencyKey)
	holdId, err := s.repository.InsertBookingHolds(ctx, params.UserId, params.IdempotencyKey, expiry)
	if err == nil {
		logger.Debug("Created new booking hold with ID", *holdId)
		err1 := s.repository.InsertBookingHoldItems(ctx, *holdId, params.Availabilities, expiry, params.ProId)
		if err1 != nil {
			logger.Error("Failed to insert booking hold items", err1)
			return nil, errors.Join(err, err1)
		}
		logger.Debug("Successfully inserted booking hold items for hold ID", *holdId)
	} else if errors.Is(err, utils.ErrIdempotencyKeyIsUsed) {
		logger.Debug("Idempotency key already used, retrieving existing hold")
		hold, err1 := s.repository.GetBookingHold(ctx, params.UserId, params.IdempotencyKey)
		if err1 != nil {
			logger.Error("Failed to get existing booking hold", err1)
			return nil, errors.Join(err, err1)
		}
		logger.Debug("Retrieved existing booking hold with ID", hold.ID)
		err1 = s.repository.EnsureAvailabilitiesBelongToProfessional(ctx, params.Availabilities, params.ProId)
		if err1 != nil {
			logger.Error("Availabilities validation failed", err1)
			return nil, errors.Join(err, err1)
		}
		logger.Debug("Validated availabilities belong to professional", params.ProId)
		holdId = &hold.ID
	} else {
		logger.Error("Failed to insert booking hold", err)
		return nil, err
	}

	return holdId, nil
}
