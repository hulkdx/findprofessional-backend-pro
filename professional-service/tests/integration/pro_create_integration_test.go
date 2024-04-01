package integration_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func CreateProTest(t *testing.T, db *sql.DB) {
	handler := router.Handler(NewTestController(db))

	t.Run("not found a record, create a new record", func(t *testing.T) {
		// Arrange
		bodyRequest := professional.CreateRequest{
			Email:     "test@gmail.com",
			Password:  "P@ssw0rd123",
			FirstName: "John",
			LastName:  "Doe",
			SkypeId:   "john_doe_skype",
			AboutMe:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		}
		request := createProRequest(bodyRequest)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusCreated)
	})
}

func createProRequest(body professional.CreateRequest) *http.Request {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	return NewJsonRequest("PUT", "/professional", &buf)
}
