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
		repository := &MockRepository{findAllSuccess: data}
		request, _ := http.NewRequest("GET", "/professionals", nil)
		response := httptest.NewRecorder()
		controller := &Controller{service: NewService(repository)}
		// Act
		controller.FindAllProfessional(response, request)
		// Assert
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, response.Body.String(), "[]")
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
		expected := []Professional{
			{
				ID:    1,
				Email: "test1@gmail.com",
			},
			{
				ID:    2,
				Email: "test2@gmail.com",
			},
		}
		repository := &MockRepository{findAllSuccess: data}
		request, _ := http.NewRequest("GET", "/professionals", nil)
		response := httptest.NewRecorder()
		controller := &Controller{service: NewService(repository)}
		// Act
		controller.FindAllProfessional(response, request)
		// Assert
		assert.Equal(t, response.Code, http.StatusOK)
		assert.EqualJSON(t, response.Body.String(), expected)
	})
}

type MockRepository struct {
	findAllSuccess []Professional
	findAllError   error
}

func (r *MockRepository) FindAll(fields ...string) ([]Professional, error) {
	// Mimic the original filtering in repository
	filter := []Professional{}
	for _, pro := range r.findAllSuccess {
		fpro := Professional{}
		for _, field := range fields {
			switch field {
			case "ID":
				fpro.ID = pro.ID
			case "Email":
				fpro.Email = pro.Email
			case "Password":
				fpro.Password = pro.Password
			case "Created_at":
				fpro.CreatedAt = pro.CreatedAt
			case "Updated_at":
				fpro.UpdatedAt = pro.UpdatedAt
			}
		}
		filter = append(filter, fpro)
	}
	return filter, r.findAllError
}
