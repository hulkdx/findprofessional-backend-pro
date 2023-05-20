package integration_test

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"gorm.io/gorm"
)

func UpdateProfessionalTest(t *testing.T, db *sql.DB, gdb *gorm.DB) {
	handler := router.Handler(NewTestController(db))

	t.Run("Empty database", func(t *testing.T) {
		// Arrange
		id := 1
		requestBody := `{ "email": "new@email.com" }`
		request := NewJsonRequest("POST", fmt.Sprintf("/professional/%d", id), strings.NewReader(requestBody))
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
		gdb.Create(record)
		defer gdb.Delete(record)

		requestBody := `{ "email": "new@email.com" }`
		request := NewJsonRequest("POST", fmt.Sprintf("/professional/%d", id), strings.NewReader(requestBody))
		response := httptest.NewRecorder()
		expected := &professional.Professional{
			ID:    id,
			Email: "new@email.com",
		}
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)
		assert.EqualJSON(t, response.Body.String(), expected)
	})
}

func NewJsonRequest(method, url string, body io.Reader) *http.Request {
	request, _ := http.NewRequest(method, url, body)
	request.Header.Set("Content-Type", "application/json")
	return request
}
