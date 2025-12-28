package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdateProfessionalTest(t *testing.T, db *pgxpool.Pool) {
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

		requestBody := `
			{
				"email": "test@gmail.com",
				"firstName": "John",
				"lastName": "Doe",
				"coachType": "Fitness Coach",
				"priceNumber": 100,
				"priceCurrency": "USD",
				"profileImageUrl": "http://example.com/images/john.jpg",
				"description": "Experienced fitness coach with 10 years of experience.",
				"sessionPlatform": "meet",
				"sessionLink": "https://meet.google.com/abc-defg-hij"
			}
		`
		request := NewJsonRequest("POST", "/professional", strings.NewReader(requestBody))
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)
	})

	t.Run("should update success when email nil but price set", func(t *testing.T) {
		// Arrange
		id := int64(1)
		userService.UserId = id
		record := &professional.Professional{
			ID:    int64(id),
			Email: "emailofidone@email.com",
		}
		d1 := insertPro(t, db, *record)
		defer d1()

		requestBody := `
			{
				"firstName": "John",
				"lastName": "Doe",
				"coachType": "Fitness Coach",
				"priceNumber": 100,
				"priceCurrency": "USD",
				"profileImageUrl": "http://example.com/images/john.jpg",
				"description": "Experienced fitness coach with 10 years of experience.",
				"sessionPlatform": "meet",
				"sessionLink": "https://meet.google.com/abc-defg-hij"
			}
		`
		request := NewJsonRequest("POST", "/professional", strings.NewReader(requestBody))
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)
	})
}
