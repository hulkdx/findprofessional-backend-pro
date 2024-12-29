package professional

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func TestGetAvailability(t *testing.T) {
	t.Run("authorize", func(t *testing.T) {
		// Arrange
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}
		// Act
		controller.GetAvailability(httptest.NewRecorder(), getAvailibilityRequest())
		// Assert
		assert.Equal(t, userService.GetAuthenticatedUserIdCalled, true)
	})
}

func getAvailibilityRequest() *http.Request {
	request, _ := http.NewRequest("GET", "/professional/availability", nil)
	return request
}
