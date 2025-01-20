package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AddAvailabilityTest(t *testing.T, db *pgxpool.Pool) {
	timeProvider := &FakeTimeProvider{}
	handler := router.Handler(NewTestControllerWithTimeProvider(db, timeProvider))

	t.Run("empty availability", func(t *testing.T) {
		// Arrange
		d1 := insertEmptyPro(t, db)
		defer d1()

		requestBody := professional.UpdateAvailabilityRequest{
			Items: []professional.UpdateAvailabilityItemRequest{
				{
					Date: "2023-01-01",
					From: "10:00:00",
					To:   "11:00:00",
				},
			},
		}
		request := NewJsonRequestBody("POST", "/professional/availability", requestBody)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)

		// Arrange
		response = httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, NewJsonRequest("GET", "/professional/availability", nil))
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)
		response_model := []professional.Availability{}
		Unmarshal(response, &response_model)
		assert.Equal(t, len(response_model), 1)
		assert.Equal(t, response_model[0].Date.String(), requestBody.Items[0].Date)
		assert.Equal(t, response_model[0].From.String(), requestBody.Items[0].From)
		assert.Equal(t, response_model[0].To.String(), requestBody.Items[0].To)
	})
}
