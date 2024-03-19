package professional

import (
	"net/http/httptest"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func TestCreateProfessional(t *testing.T) {
	t.Run("authorize", func(t *testing.T) {
		// Arrange
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}
		// Act
		controller.Create(httptest.NewRecorder(), findAllRequest())
		// Assert
		assert.Equal(t, userService.IsAuthenticatedCalled, true)
	})
}
