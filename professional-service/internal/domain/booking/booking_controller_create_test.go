package booking

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/mocks"
)

func TestCreateBooking(t *testing.T) {
	t.Run("authorize", func(t *testing.T) {
		// Arrange
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := NewController(userService, NewService(FakeRepository{}, &mocks.FakePaymentService{}))
		// Act
		controller.Create(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/booking", nil))
		// Assert
		assert.Equal(t, userService.GetAuthenticatedUserIdCalled, true)
	})
}

// --- Test doubles ---
type MockUserServiceAlwaysAuthenticated struct {
	IsAuthenticatedCalled        bool
	GetAuthenticatedUserIdCalled bool
}

func (m *MockUserServiceAlwaysAuthenticated) IsAuthenticated(context.Context, string) bool {
	m.IsAuthenticatedCalled = true
	return true
}

func (m *MockUserServiceAlwaysAuthenticated) Login(ctx context.Context, email string, password string) (string, error) {
	return "", nil
}

func (m *MockUserServiceAlwaysAuthenticated) GetAuthenticatedUserId(ctx context.Context, auth string) (int64, error) {
	m.GetAuthenticatedUserIdCalled = true
	return -2, nil
}

type FakeRepository struct {
}

func (r FakeRepository) GetPriceAndCurrency(ctx context.Context, proId string) (int64, string, error) {
	return 5000, "eur", nil
}

func (r FakeRepository) InsertBooking(ctx context.Context, userId int64, proId string, req *booking_model.CreateBookingRequest) (int64, error) {
	return 1, nil
}
