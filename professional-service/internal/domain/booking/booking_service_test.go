package booking

import (
	"context"
	"testing"
	"time"

	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/mocks"
)

func TestCreateBookingService(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repository := &FakeRepository{}
	//timProvider := utils.TimeProvider()
	service := NewService(repository, &mocks.FakePaymentService{}, nil)
	params := &CreateParams{}
	// Act
	_, err := service.create(ctx, params)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
}

// --- Test doubles ---

type FakeRepository struct {
}

func (r *FakeRepository) WithTx(ctx context.Context, fn WithTxFunc) (*booking_model.CreateBookingResponse, error) {
	return nil, nil
}

func (r *FakeRepository) InsertBookingHolds(ctx context.Context, UserId int64, IdempotencyKey string, expiry time.Time) (*int64, error) {
	return nil, nil
}
