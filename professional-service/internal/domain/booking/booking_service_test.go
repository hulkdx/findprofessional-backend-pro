package booking

import (
	"context"
	"testing"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
)

func TestCreateBookingService(t *testing.T) {
	/*
		t.Run("success", func(t *testing.T) {
			// Arrange
			repository := &FakeRepository{}
			timProvider := &mocks.FakeTimeProvider{}
			service := NewService(repository, &mocks.FakePaymentService{}, timProvider)
			params := &CreateParams{}
			// Act
			_, err := service.create(context.Background(), params)
			// Assert
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, repository.IsInsertBookingHoldsCalled, true)
		})

		t.Run("If InsertBookingHolds conflicted, then it should return a error", func(t *testing.T) {
			// Arrange
			repository := &FakeRepository{}
			repository.InsertBookingHoldsResponse = nil
			repository.InsertBookingHoldsErrResponse = utils.ErrIdempotencyKeyIsUsed
			timProvider := &mocks.FakeTimeProvider{}
			service := NewService(repository, &mocks.FakePaymentService{}, timProvider)
			params := &CreateParams{}
			// Act
			_, err := service.create(context.Background(), params)
			// Assert
			assert.Equal(t, err, utils.ErrIdempotencyKeyIsUsed)
		})
	*/
}

// --- Test doubles ---

type FakeRepository struct {
	IsInsertBookingHoldsCalled    bool
	InsertBookingHoldsResponse    *int64
	InsertBookingHoldsErrResponse error
}

func (r *FakeRepository) WithTx(ctx context.Context, fn WithTxFunc) (*bookingmodel.CreateBookingResponse, error) {
	return nil, nil
}

func (r *FakeRepository) InsertBookingHolds(ctx context.Context, UserId int64, IdempotencyKey string, expiry time.Time) (*int64, error) {
	r.IsInsertBookingHoldsCalled = true
	return r.InsertBookingHoldsResponse, r.InsertBookingHoldsErrResponse
}

func (r *FakeRepository) GetBookingHold(ctx context.Context, userId int64, idempotencyKey string) (*bookingmodel.BookingHold, error) {
	return nil, nil
}
