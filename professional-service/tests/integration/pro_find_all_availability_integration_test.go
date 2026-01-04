package integration_test

import (
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

func FindAllAvailabilityProfessionalTest(t *testing.T, db *pgxpool.Pool) {
	timeProvider := &mocks.FakeTimeProvider{}
	handler := router.Handler(NewTestControllerWithTimeProvider(db, timeProvider))

	t.Run("empty availability", func(t *testing.T) {
		// Arrange
		expected := []professional.Availability{}
		d1 := insertEmptyPro(t, db)
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
		timeProvider.NowTime = time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

		expected := []professional.Availability{
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
		ids, d2 := insertAvailability(t, db, expected...)
		expected[0].ID = ids[0]
		expected[1].ID = ids[1]

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

	t.Run("should not show availabilities that are older than current time", func(t *testing.T) {
		// Arrange
		timeProvider.NowTime = time.Date(2025, 9, 24, 12, 12, 0, 0, time.UTC)

		records := []professional.Availability{
			{
				From: time.Date(2023, 11, 4, 5, 30, 0, 0, time.UTC),
				To:   time.Date(2023, 11, 4, 6, 30, 0, 0, time.UTC),
			},
			{
				From: time.Date(2020, 11, 4, 15, 30, 0, 0, time.UTC),
				To:   time.Date(2020, 11, 4, 16, 00, 0, 0, time.UTC),
			},
			{
				From: time.Date(2025, 9, 23, 5, 30, 0, 0, time.UTC),
				To:   time.Date(2025, 9, 23, 6, 30, 0, 0, time.UTC),
			},
			{
				From: time.Date(2025, 9, 24, 11, 30, 0, 0, time.UTC),
				To:   time.Date(2025, 9, 24, 12, 00, 0, 0, time.UTC),
			},
		}
		d1 := insertEmptyPro(t, db)
		_, d2 := insertAvailability(t, db, records...)
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
		assert.Equal(t, len(response_model[0].Availability), 0)
	})

	t.Run("should not show availabilities that are booked", func(t *testing.T) {
		// Arrange
		timeProvider.NowTime = time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

		userId := int64(2)
		proId := int64(1)

		pro := professional.Professional{
			ID:            proId,
			PriceNumber:   Int(0),
			PriceCurrency: String(""),
			Pending:       false,
		}

		records := []professional.Availability{
			{
				ProfessionalID: proId,
				From:           time.Date(2023, 11, 4, 5, 30, 0, 0, time.UTC),
				To:             time.Date(2023, 11, 4, 6, 30, 0, 0, time.UTC),
			},
		}
		d1 := insertPro(t, db, pro)
		defer d1()
		avIds, d2 := insertAvailability(t, db, records...)
		defer d2()
		d3 := insertUserWithId(t, db, userId)
		defer d3()
		bookingId, d4 := insertBooking(t, db, userId, proId, "pending", "", "abc")
		defer d4()
		d5 := insertBookingItems(t, db, TestBookingItems{BookingID: bookingId, AvailabilityID: avIds[0]})
		defer d5()

		request := NewJsonRequest("GET", "/professional", nil)
		response := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(response, request)
		// Assert
		assert.Equal(t, response.Code, http.StatusOK)

		responseModel := []professional.Professional{}
		Unmarshal(response, &responseModel)

		assert.Equal(t, len(responseModel), 1)
		assert.Equal(t, len(responseModel[0].Availability), 0)
	})
}
