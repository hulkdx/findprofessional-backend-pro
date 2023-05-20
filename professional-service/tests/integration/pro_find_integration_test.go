package integration_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"gorm.io/gorm"
)

func FindProfessionalTest(t *testing.T, db *sql.DB, gdb *gorm.DB) {
	handler := router.Handler(NewTestController(db))

	t.Run("Empty database", func(t *testing.T) {
		// Arrange
		id := 1
		request := NewJsonRequest("GET", fmt.Sprintf("/professional/%d", id), nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	t.Run("found a record", func(t *testing.T) {
		// Arrange
		id := 1
		record := &professional.Professional{
			ID:    id,
			Email: "emailofidone@email.com",
		}
		request := NewJsonRequest("GET", fmt.Sprintf("/professional/%d", id), nil)
		response := httptest.NewRecorder()
		gdb.Create(record)
		defer gdb.Delete(record)
		expected := &professional.Professional{
			ID:    id,
			Email: "emailofidone@email.com",
		}
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)
		assert.EqualJSON(t, response.Body.String(), expected)
	})
}
