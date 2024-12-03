package integration_test

import (
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func UpdateProfessionalTest(t *testing.T, db *sql.DB) {
	userService := MockUserService{}
	handler := router.Handler(NewTestControllerWithUserService(db, &userService))

	t.Run("Empty database", func(t *testing.T) {
		// Arrange
		userService.UserId = 1
		requestBody := `{ "email": "new@email.com" }`
		request := NewJsonRequest("POST", "/professional", strings.NewReader(requestBody))
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	t.Run("found a record", func(t *testing.T) {
		// Arrange
		id := int64(1)
		userService.UserId = id
		record := &professional.Professional{
			ID:    int64(id),
			Email: "emailofidone@email.com",
		}
		d1 := insertPro(t, db, *record)
		defer d1()

		requestBody := `{ "email": "new@email.com" }`
		request := NewJsonRequest("POST", "/professional", strings.NewReader(requestBody))
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)
	})
}

func NewJsonRequest(method, url string, body io.Reader) *http.Request {
	request, _ := http.NewRequest(method, url, body)
	request.Header.Set("Content-Type", "application/json")
	return request
}
