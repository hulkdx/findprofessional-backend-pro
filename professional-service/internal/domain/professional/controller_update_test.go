package professional

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func TestUpdateProfessional(t *testing.T) {
	t.Run("authorize", func(t *testing.T) {
		// Arrange
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}
		// Act
		controller.Update(httptest.NewRecorder(), updateRequest(1, ""))
		// Assert
		assert.Equal(t, userService.IsAuthenticatedCalled, true)
	})

	t.Run("not found the id", func(t *testing.T) {
		// Arrange
		id := 1
		requestBody := `{ "email": "new@email.com" }`
		request := updateRequest(id, requestBody)
		response := httptest.NewRecorder()
		controller := updateController(sql.ErrNoRows)
		// Act
		controller.Update(response, request)
		// Assert
		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	t.Run("found in repository", func(t *testing.T) {
		// Arrange
		id := 1
		requestBody := `{ "email": "new@email.com" }`
		request := updateRequest(id, requestBody)
		response := httptest.NewRecorder()
		controller := updateController(nil)
		// Act
		controller.Update(response, request)
		// Assert
		assert.Equal(t, response.Code, http.StatusOK)
	})

	// TODO: add test for invalid email address
}

func updateController(updateError error) *Controller {
	repository := &FakeRepository{updateError: updateError}
	return &Controller{
		service:     NewService(repository),
		userService: &MockUserServiceAlwaysAuthenticated{},
	}
}

func updateRequest(id int, body string) *http.Request {
	request, _ := http.NewRequest("POST", fmt.Sprintf("/professional/%d", id), nil)
	return request
}
