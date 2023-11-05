package integration_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/civil"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func FindAllAvailabilityProfessionalTest(t *testing.T, db *sql.DB) {
	handler := router.Handler(NewTestController(db))
	t.Run("empty availability", func(t *testing.T) {
		// Arrange
		expected := []professional.Availability{}
		d1 := insertEmptyPro(db)
		defer d1()

		request := NewJsonRequest("GET", "/professional", nil)
		response := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(response, request)

		// Assert
		assert.Equal(t, response.Code, http.StatusOK)

		response_model := []professional.Professional{}
		Unmarshal(response, &response_model)

		assert.Equal(t, len(response_model), 1)
		assert.EqualAnyOrder(t, response_model[0].Availability, expected)
	})

	t.Run("some availabilities", func(t *testing.T) {
		// Arrange
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

		d1 := insertEmptyPro(db)
		d2 := insertAvailability(db, expected...)
		defer d2()
		defer d1()

		request := NewJsonRequest("GET", "/professional", nil)
		response := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(response, request)

		// Assert
		assert.Equal(t, response.Code, http.StatusOK)

		response_model := []professional.Professional{}
		Unmarshal(response, &response_model)

		assert.Equal(t, len(response_model), 1)
		assert.EqualAnyOrder(t, response_model[0].Availability, expected)
	})
}
