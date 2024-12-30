package integration_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func GetAvailabilityTest(t *testing.T, db *sql.DB) {
	timeProvider := &FakeTimeProvider{}
	handler := router.Handler(NewTestControllerWithTimeProvider(db, timeProvider))

	t.Run("empty availability", func(t *testing.T) {
		// Arrange
		request := NewJsonRequest("GET", "/professional/availability", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)

		response_model := []professional.Availability{}
		Unmarshal(response, &response_model)
		assert.Equal(t, len(response_model), 0)
	})

	t.Run("some availability", func(t *testing.T) {
		// Arrange
		timeProvider.NowTime = time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

		expected := []professional.Availability{
			{
				ID:   0,
				Date: civil.Date{Year: 2023, Month: 11, Day: 4},
				From: civil.Time{Hour: 5, Minute: 30},
				To:   civil.Time{Hour: 6, Minute: 30},
			},
			{
				ID:   0,
				Date: civil.Date{Year: 2020, Month: 11, Day: 4},
				From: civil.Time{Hour: 15, Minute: 30},
				To:   civil.Time{Hour: 16, Minute: 00},
			},
		}

		d1 := insertEmptyPro(t, db)
		d2 := insertAvailability(t, db, expected...)
		defer d2()
		defer d1()

		request := NewJsonRequest("GET", "/professional/availability", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)

		response_model := []professional.Availability{}
		Unmarshal(response, &response_model)
		assert.Equal(t, response_model, expected)
	})
}
