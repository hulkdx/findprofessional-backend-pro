package booking

import (
	"context"
	"errors"
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

	// TODO: Insert booking into the database as holding state (next step)
	// TODO: Call payment service to create a payment intent (later)
	return nil, nil
}

func (s *BookingService) validate(ctx context.Context, userId int64, proId string, req CreateBookingRequest) error {
	priceNumber, currency, err := s.repository.GetPriceAndCurrency(ctx, proId)
	if err != nil {
		return err
	}
	amountsInCents := priceNumber * int64(len(req.Slots))
	if amountsInCents != req.AmountInCents {
		return errors.New("invalid amount_in_cents")
	}
	if currency != req.Currency {
		return errors.New("invalid currency")
	}
	return nil
}
