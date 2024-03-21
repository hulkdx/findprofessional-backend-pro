package professional

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

	t.Run("valid emails", func(t *testing.T) {
		validEmails := []string{
			"new@email.com",
			"50charemailxxxxxxxxxxxxxxxxxxxxxxxxxxxxx@gmail.com",
		}
		for _, email := range validEmails {
			// Arrange
			id := 1
			requestBody := fmt.Sprintf(`{ "email": "%s" }`, email)
			request := updateRequest(id, requestBody)
			response := httptest.NewRecorder()
			controller := updateController(nil)
			// Act
			controller.Update(response, request)
			// Assert
			assert.Equal(t, response.Code, http.StatusOK)
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		invalidEmails := []string{
			"",
			"23123",
			"space email@gmail.com",
			"51charemailxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx@gmail.com",
		}
		for _, email := range invalidEmails {
			// Arrange
			id := 1
			requestBody := fmt.Sprintf(`{ "email": "%s" }`, email)
			request := updateRequest(id, requestBody)
			response := httptest.NewRecorder()
			controller := updateController(nil)
			// Act
			controller.Update(response, request)
			// Assert
			assert.Equal(t, response.Code, http.StatusBadRequest)
		}
	})
}

func updateController(updateError error) *Controller {
	repository := &FakeRepository{updateError: updateError}
	return &Controller{
		service:     NewService(repository),
		userService: &MockUserServiceAlwaysAuthenticated{},
	}
}

func updateRequest(id int, body string) *http.Request {
	request, _ := http.NewRequest("POST", fmt.Sprintf("/professional/%d", 400), strings.NewReader(body))
	return request
}
