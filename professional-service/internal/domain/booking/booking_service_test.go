package booking

import (
	"context"
	"testing"
	"time"

	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/mocks"
)

func TestCreateBookingService(t *testing.T) {
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
		repository.InsertBookingHoldsErrResponse = utils.ErrIdempotencyKeyExpired
		timProvider := &mocks.FakeTimeProvider{}
		service := NewService(repository, &mocks.FakePaymentService{}, timProvider)
		params := &CreateParams{}
		// Act
		_, err := service.create(context.Background(), params)
		// Assert
		assert.Equal(t, err, utils.ErrIdempotencyKeyExpired)
	})
}

// --- Test doubles ---

type FakeRepository struct {
	IsInsertBookingHoldsCalled    bool
	InsertBookingHoldsResponse    *int64
	InsertBookingHoldsErrResponse error
}

func (r *FakeRepository) WithTx(ctx context.Context, fn WithTxFunc) (*booking_model.CreateBookingResponse, error) {
	return nil, nil
}

func (r *FakeRepository) InsertBookingHolds(ctx context.Context, UserId int64, IdempotencyKey string, expiry time.Time) (*int64, error) {
	r.IsInsertBookingHoldsCalled = true
	return r.InsertBookingHoldsResponse, r.InsertBookingHoldsErrResponse
}
