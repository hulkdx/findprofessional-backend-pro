package integration_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"gorm.io/gorm"
)

func FindAllProfessionalTest(t *testing.T, db *sql.DB, gdb *gorm.DB) {
	handler := router.Handler(NewTestController(db))

	t.Run("Empty professionals", func(t *testing.T) {
		// Arrange
		request := NewJsonRequest("GET", "/professional", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)
		assert.EqualJSON(t, response.Body.String(), []string{})
	})

	t.Run("some professional, only show valid data", func(t *testing.T) {
		// Arrange
		now := time.Now()
		records := []professional.Professional{
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
		gdb.Create(records)
		defer gdb.Delete(records)
		expected := []professional.Professional{
			{
				ID:    1,
				Email: "test1@gmail.com",
			},
			{
				ID:    2,
				Email: "test2@gmail.com",
			},
		}
		request := NewJsonRequest("GET", "/professional", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Assert
		assert.Equal(t, response.Code, http.StatusOK)
		assert.EqualJSON(t, response.Body.String(), expected)
	})
}

func NewTestController(db *sql.DB) *professional.Controller {
	controller := professional.NewController(
		professional.NewService(professional.NewRepository(db)),
		&MockUserService{},
	)
	return controller
}

type MockUserService struct{}

func (m *MockUserService) IsAuthenticated(ctx context.Context, auth string) bool {
	return true
}
