package booking

import (
	"context"

	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/payment"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/logger"
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
	err := s.validate(ctx, proId, req)
	if err != nil {
		logger.Error("validation error: %v", err)
		return nil, err
	}

	bookingId, err := s.repository.InsertBooking(ctx, userId, proId, req)
	if err != nil {
		logger.Error("repository InsertBooking error:", err)
		return nil, err
	}

	paymentIntentResponse, err := s.paymentService.CreatePaymentIntent(ctx, userId, req, auth)
	if err != nil {
		logger.Error("paymentService CreatePaymentIntent error:", err)
		return nil, err
	}
	return &booking_model.CreateBookingResponse{
		BookingID:             bookingId,
		PaymentIntentResponse: *paymentIntentResponse,
	}, nil
}

func (s *BookingService) validate(ctx context.Context, proId string, req *booking_model.CreateBookingRequest) error {
	slotSize := len(req.Slots)
	if slotSize == 0 || slotSize > 10 {
		return utils.ErrInvalidSlotSize
	}

	priceNumber, currency, err := s.repository.GetPriceAndCurrency(ctx, proId)
	if err != nil {
		return utils.ErrValidationDatabase
	}
	amountsInCents := priceNumber * int64(len(req.Slots))
	if amountsInCents != req.AmountInCents {
		return utils.ErrAmountInCentsMismatch
	}
	if currency != req.Currency {
		return utils.ErrCurrencyMismatch
	}
	return nil
}
