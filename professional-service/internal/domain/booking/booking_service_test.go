package booking

import (
	"context"
	"testing"

	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/mocks"
)

func TestCreateBookingService(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repository := &FakeRepository{}
	service := NewService(repository, &mocks.FakePaymentService{})
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

func (r *FakeRepository) WithTx(ctx context.Context, fn Func) (*booking_model.CreateBookingResponse, error) {
	return nil, nil
}
