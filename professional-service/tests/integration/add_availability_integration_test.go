package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/civil"
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
		defer db.Exec(context.Background(), `DELETE FROM professional_availability`)

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

	t.Run("duplicate availability, update them", func(t *testing.T) {
		// Arrange
		d1 := insertEmptyPro(t, db)
		databaseAvailability := []professional.Availability{
			{
				ID:   0,
				Date: civil.Date{Year: 2023, Month: 1, Day: 1},
				From: civil.Time{Hour: 5, Minute: 30},
				To:   civil.Time{Hour: 6, Minute: 30},
			},
		}
		d2 := insertAvailability(t, db, databaseAvailability...)

		defer d1()
		defer d2()
		defer db.Exec(context.Background(), `DELETE FROM professional_availability`)

		newAvailability := professional.UpdateAvailabilityRequest{
			Items: []professional.UpdateAvailabilityItemRequest{
				{
					Date: "2023-01-01",
					From: "04:00:00",
					To:   "07:30:00",
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
		assert.Equal(t, response_model[0].Date.String(), newAvailability.Items[0].Date)
		assert.Equal(t, response_model[0].From.String(), newAvailability.Items[0].From)
		assert.Equal(t, response_model[0].To.String(), newAvailability.Items[0].To)
	})

	t.Run("non-duplicate availability, should be returned", func(t *testing.T) {
		// Arrange
		d1 := insertEmptyPro(t, db)
		databaseAvailability := []professional.Availability{
			{
				ID:   0,
				Date: civil.Date{Year: 2023, Month: 1, Day: 2},
				From: civil.Time{Hour: 5, Minute: 30},
				To:   civil.Time{Hour: 6, Minute: 30},
			},
		}
		d2 := insertAvailability(t, db, databaseAvailability...)

		defer d1()
		defer d2()
		defer db.Exec(context.Background(), `DELETE FROM professional_availability`)

		newAvailability := professional.UpdateAvailabilityRequest{
			Items: []professional.UpdateAvailabilityItemRequest{
				{
					Date: "2023-01-01",
					From: "04:00:00",
					To:   "07:30:00",
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

	t.Run("same-date availability, should be removed", func(t *testing.T) {
		// Arrange
		d1 := insertEmptyPro(t, db)
		databaseAvailability := []professional.Availability{
			{
				ID:   0,
				Date: civil.Date{Year: 2023, Month: 1, Day: 1},
				From: civil.Time{Hour: 8, Minute: 30},
				To:   civil.Time{Hour: 9, Minute: 30},
			},
		}
		d2 := insertAvailability(t, db, databaseAvailability...)

		defer d1()
		defer d2()
		defer db.Exec(context.Background(), `DELETE FROM professional_availability`)

		newAvailability := professional.UpdateAvailabilityRequest{
			Items: []professional.UpdateAvailabilityItemRequest{
				{
					Date: "2023-01-01",
					From: "04:00:00",
					To:   "07:30:00",
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
		assert.Equal(t, response_model[0].Date.String(), newAvailability.Items[0].Date)
		assert.Equal(t, response_model[0].From.String(), newAvailability.Items[0].From)
		assert.Equal(t, response_model[0].To.String(), newAvailability.Items[0].To)
	})
}
