package integration_test

import (
	"database/sql"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func ListProfessionalTest(t *testing.T, db *sql.DB, gdb *gorm.DB) {
	sut := router.Handler(db)

	t.Run("Empty professionals", func(t *testing.T) {
		// Arrange
		request, _ := http.NewRequest("GET", "/professionals", nil)
		response := httptest.NewRecorder()
		// Act
		sut.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)
		assert.EqualJSON(t, response.Body.String(), []string{})
	})
	t.Run("some professional, only show valid data", func(t *testing.T) {
		// Arrange
		now := time.Now()
		data := []professional.Professional{
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
		gdb.Create(data)
		defer gdb.Delete(data)
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
		request, _ := http.NewRequest("GET", "/professionals", nil)
		response := httptest.NewRecorder()
		// Act
		sut.ServeHTTP(response, request)
		// Assert
		assert.Equal(t, response.Code, http.StatusOK)
		assert.EqualJSON(t, response.Body.String(), expected)
	})
}
