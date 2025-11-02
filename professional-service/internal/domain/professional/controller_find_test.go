package professional

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func TestFindProfessional(t *testing.T) {
	t.Run("authorize", func(t *testing.T) {
		// Arrange
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}
		// Act
		controller.Find(httptest.NewRecorder(), findRequest(1))
		// Assert
		assert.Equal(t, userService.IsAuthenticatedCalled, true)
	})
	t.Run("empty repository", func(t *testing.T) {
		// Arrange
		id := 1
		response := httptest.NewRecorder()
		controller := findController(nil, sql.ErrNoRows)
		// Act
		controller.Find(response, findRequest(id))
		// Assert
		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	t.Run("found an id", func(t *testing.T) {
		// Arrange
		id := 1
		record := &model_professional.Professional{
			ID:    1,
			Email: "emailofidone@email.com",
		}
		response := httptest.NewRecorder()
		controller := findController(record, nil)
		// Act
		controller.Find(response, findRequest(id))
		// Assert
		assert.Equal(t, response.Code, http.StatusOK)
		assert.EqualJSON(t, response.Body.String(), record)
	})
}

func findController(findByIdSuccess *model_professional.Professional, findByIdError error) *Controller {
	repository := &FakeRepository{}
	if findByIdSuccess != nil {
		repository.findByIdSuccess = *findByIdSuccess
	}
	if findByIdError != nil {
		repository.findByIdError = findByIdError
	}
	return &Controller{
		service:     NewService(repository),
		userService: &MockUserServiceAlwaysAuthenticated{},
	}
}

func findRequest(id int) *http.Request {
	request, _ := http.NewRequest("GET", fmt.Sprintf("/professional/%d", id), nil)
	return request
}
