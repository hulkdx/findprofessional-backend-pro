package booking

import (
	"context"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

type BookingService struct {
	repository Repository
}

func NewService(repository Repository) *BookingService {
	return &BookingService{
		repository: repository,
	}
}

func (s *BookingService) Create(ctx context.Context, userId int64, proId string, req CreateBookingRequest) (*CreateBookingResponse, error) {
	err := s.validate(ctx, userId, proId, req)
	if err != nil {
		return nil, err
	}

	err = s.repository.InsertBooking(ctx, userId, proId, req)
	if err != nil {
		return nil, err
	}
	// TODO: Call payment service to create a payment intent (later)
	return nil, nil
}

func (s *BookingService) validate(ctx context.Context, userId int64, proId string, req CreateBookingRequest) error {
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
