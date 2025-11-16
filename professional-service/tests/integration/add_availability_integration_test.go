package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/mocks"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AddAvailabilityTest(t *testing.T, db *pgxpool.Pool) {
	timeProvider := &mocks.FakeTimeProvider{}
	handler := router.Handler(NewTestControllerWithTimeProvider(db, timeProvider))

	t.Run("empty availability", func(t *testing.T) {
		// Arrange
		d1 := insertEmptyPro(t, db)
		defer d1()
		defer db.Exec(context.Background(), `DELETE FROM professional_availability`)

		requestBody := professional.UpdateAvailabilityRequest{
			Items: []professional.UpdateAvailabilityItemRequest{
				{
					From: time.Date(2023, 1, 1, 10, 00, 0, 0, time.UTC),
					To:   time.Date(2023, 1, 1, 11, 00, 0, 0, time.UTC),
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
		assert.Equal(t, response_model[0].From, requestBody.Items[0].From)
		assert.Equal(t, response_model[0].To, requestBody.Items[0].To)
	})

	t.Run("duplicate availability, update them", func(t *testing.T) {
		// Arrange
		d1 := insertEmptyPro(t, db)
		databaseAvailability := []professional.Availability{
			{
				ID:   0,
				From: time.Date(2023, 1, 1, 5, 30, 0, 0, time.UTC),
				To:   time.Date(2023, 1, 1, 6, 30, 0, 0, time.UTC),
			},
		}
		_, d2 := insertAvailability(t, db, databaseAvailability...)

		defer d1()
		defer d2()
		defer db.Exec(context.Background(), `DELETE FROM professional_availability`)

		newAvailability := professional.UpdateAvailabilityRequest{
			Items: []professional.UpdateAvailabilityItemRequest{
				{
					From: time.Date(2023, 1, 1, 4, 00, 0, 0, time.UTC),
					To:   time.Date(2023, 1, 1, 7, 30, 0, 0, time.UTC),
				},
			},
		}
		request := NewJsonRequestBody("POST", "/professional/availability", newAvailability)
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
		assert.Equal(t, response_model[0].From, newAvailability.Items[0].From)
		assert.Equal(t, response_model[0].To, newAvailability.Items[0].To)
	})

	t.Run("non-duplicate availability, should be returned", func(t *testing.T) {
		// Arrange
		d1 := insertEmptyPro(t, db)
		databaseAvailability := []professional.Availability{
			{
				ID:   0,
				From: time.Date(2023, 1, 2, 5, 30, 0, 0, time.UTC),
				To:   time.Date(2023, 1, 2, 6, 30, 0, 0, time.UTC),
			},
		}
		_, d2 := insertAvailability(t, db, databaseAvailability...)

		defer d1()
		defer d2()
		defer db.Exec(context.Background(), `DELETE FROM professional_availability`)

		newAvailability := professional.UpdateAvailabilityRequest{
			Items: []professional.UpdateAvailabilityItemRequest{
				{
					From: time.Date(2023, 1, 1, 4, 00, 0, 0, time.UTC),
					To:   time.Date(2023, 1, 1, 7, 30, 0, 0, time.UTC),
				},
			},
		}
		request := NewJsonRequestBody("POST", "/professional/availability", newAvailability)
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
		assert.Equal(t, len(response_model), 2)
	})

	t.Run("on empty request body response should fail with", func(t *testing.T) {
		// Arrange
		d1 := insertEmptyPro(t, db)
		defer d1()
		defer db.Exec(context.Background(), `DELETE FROM professional_availability`)

		request := NewJsonRequestBody("POST", "/professional/availability", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusBadRequest)
	})
}
