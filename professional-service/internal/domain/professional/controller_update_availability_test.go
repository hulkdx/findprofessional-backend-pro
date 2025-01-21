package professional

import (
	"net/http/httptest"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func TestUpdateAvailability(t *testing.T) {
	t.Run("authorize", func(t *testing.T) {
		// Arrange
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}
		// Act
		controller.UpdateAvailability(httptest.NewRecorder(), getAvailibilityRequest())
		// Assert
		assert.Equal(t, userService.GetAuthenticatedUserIdCalled, true)
	})
}
