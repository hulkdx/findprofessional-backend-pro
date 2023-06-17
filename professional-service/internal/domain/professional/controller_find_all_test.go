package professional

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func TestFindAllProfessional(t *testing.T) {
	t.Run("authorize", func(t *testing.T) {
		// Arrange
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}
		// Act
		controller.FindAll(httptest.NewRecorder(), findAllRequest())
		// Assert
		assert.Equal(t, userService.IsAuthenticatedCalled, true)
	})
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
	t.Run("some professional", func(t *testing.T) {
		// Arrange
		data := []Professional{
			{
				ID:    1,
				Email: "test1@gmail.com",
			},
			{
				ID:    2,
				Email: "test2@gmail.com",
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
	t.Run("some professional with all optional fields", func(t *testing.T) {
		// Arrange
		data := []Professional{
			{
				ID:              1,
				Email:           "test1@gmail.com",
				FirstName:       "Test 1",
				LastName:        "Last 1",
				CoachType:       "Type 1",
				PriceNumber:     100,
				PriceCurrency:   "EUR",
				ProfileImageUrl: "Url 1",
			},
			{
				ID:              2,
				Email:           "test2@gmail.com",
				FirstName:       "Test 2",
				LastName:        "Last 2",
				CoachType:       "Type 2",
				PriceNumber:     100,
				PriceCurrency:   "EUR",
				ProfileImageUrl: "Url 2",
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
