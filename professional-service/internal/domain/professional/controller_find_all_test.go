package professional

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func TestFindAllProfessional(t *testing.T) {
	t.Run("empty repository", func(t *testing.T) {
		// Arrange
		data := []Professional{}
		response := httptest.NewRecorder()
		controller := findAllController(data)
		// Act
		controller.FindAll(response, findAllRequest())
		// Assert
		assert.Equal(t, response.Code, http.StatusOK)
		assert.EqualJSON(t, response.Body.String(), []string{})
	})
	t.Run("some professional, only show valid data", func(t *testing.T) {
		// Arrange
		now := time.Now()
		data := []Professional{
			{
				ID:        1,
				Email:     "test1@gmail.com",
				Password:  "some_hex_value2",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			{
				ID:        2,
				Email:     "test2@gmail.com",
				Password:  "some_hex_value2",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
		}
		response := httptest.NewRecorder()
		controller := findAllController(data)
		// Act
		controller.FindAll(response, findAllRequest())
		// Assert
		assert.Equal(t, response.Code, http.StatusOK)
		assert.EqualJSON(t, response.Body.String(), data)
	})
}

func findAllController(findAllSuccess []Professional) *Controller {
	repository := &FakeRepository{findAllSuccess: findAllSuccess}
	return &Controller{
		service:     NewService(repository),
		userService: &MockUserServiceAlwaysAuthenticated{},
	}
}

func findAllRequest() *http.Request {
	request, _ := http.NewRequest("GET", "/professionals", nil)
	return request
}
