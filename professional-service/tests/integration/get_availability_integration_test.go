package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAvailabilityTest(t *testing.T, db *pgxpool.Pool) {
	timeProvider := &FakeTimeProvider{}
	handler := router.Handler(NewTestControllerWithTimeProvider(db, timeProvider), nil)

	t.Run("empty availability", func(t *testing.T) {
		// Arrange
		request := NewJsonRequest("GET", "/professional/availability", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)

		response_model := []model_professional.Availability{}
		Unmarshal(response, &response_model)
		assert.Equal(t, len(response_model), 0)
	})

	t.Run("some availability", func(t *testing.T) {
		// Arrange
		timeProvider.NowTime = time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

		expected := []model_professional.Availability{
			{
				ID:   0,
				From: time.Date(2023, 11, 4, 5, 30, 0, 0, time.UTC),
				To:   time.Date(2023, 11, 4, 6, 30, 0, 0, time.UTC),
			},
			{
				ID:   0,
				From: time.Date(2020, 11, 4, 15, 30, 0, 0, time.UTC),
				To:   time.Date(2020, 11, 4, 16, 00, 0, 0, time.UTC),
			},
		}

		d1 := insertEmptyPro(t, db)
		_, d2 := insertAvailability(t, db, expected...)
		defer d2()
		defer d1()

		request := NewJsonRequest("GET", "/professional/availability", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)

		response_model := []model_professional.Availability{}
		Unmarshal(response, &response_model)
		assert.Equal(t, len(response_model), 2)
		assert.Equal(t, response_model[0].From, expected[0].From)
		assert.Equal(t, response_model[0].To, expected[0].To)
		assert.Equal(t, response_model[1].From, expected[1].From)
		assert.Equal(t, response_model[1].To, expected[1].To)
	})
}
