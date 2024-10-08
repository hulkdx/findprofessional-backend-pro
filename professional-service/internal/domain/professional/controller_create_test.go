package professional

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func TestCreateProfessional(t *testing.T) {
	t.Run("valid request 201", func(t *testing.T) {
		// Arrange
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: &MockUserServiceAlwaysAuthenticated{},
		}
		request := CreateRequest{
			Email:     "test@gmail.com",
			Password:  "",
			FirstName: "",
			LastName:  "",
			SkypeId:   "",
			AboutMe:   "",
		}
		response := httptest.NewRecorder()
		// Act
		controller.Create(response, createRequest(&request))
		// Assert
		assert.Equal(t, response.Code, http.StatusCreated)
	})
}

func createRequest(body *CreateRequest) *http.Request {
	jsonData, _ := json.Marshal(body)
	request, _ := http.NewRequest("PUT", "/professional", strings.NewReader(string(jsonData)))
	return request
}
