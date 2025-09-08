package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BookingCreateTest(t *testing.T, db *pgxpool.Pool) {
	timeProvider := &FakeTimeProvider{}
	handler := router.Handler(NewTestController(db), NewTestBookingController(db, timeProvider))

	t.Run("valid request", func(t *testing.T) {
		// Arrange
		records := []professional.Professional{
			{
				ID:            1,
				Email:         "test@gmail.com",
				Password:      "password",
				PriceNumber:   Int(1000),
				PriceCurrency: String("USD"),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}
		d1 := insertPro(t, db, records...)
		defer d1()

		request := NewJsonRequestBody("POST", "/professional/1/booking", booking.CreateBookingRequest{
			Slots: []booking.Slot{
				{Date: "2023-01-01", From: "10:00:00", To: "11:00:00"},
			},
			IdempotencyKey: "test-key",
			AmountInCents:  1000,
			Currency:       "USD",
		})
		request.Header.Set("Authorization", "Bearer test-token")
		response := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(response, request)

		// Assert
		assert.Equal(t, response.Code, http.StatusOK)
	})

	t.Run("invalid request - 1000 amount_in_cents with 2 slots should be 2000", func(t *testing.T) {
		// Arrange
		records := []professional.Professional{
			{
				ID:            1,
				Email:         "test@gmail.com",
				Password:      "password",
				PriceNumber:   Int(1000),
				PriceCurrency: String("USD"),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}
		d1 := insertPro(t, db, records...)
		defer d1()

		request := NewJsonRequestBody("POST", "/professional/1/booking", booking.CreateBookingRequest{
			Slots: []booking.Slot{
				{Date: "2023-01-01", From: "10:00:00", To: "11:00:00"},
				{Date: "2023-01-01", From: "12:00:00", To: "13:00:00"},
			},
			IdempotencyKey: "test-key",
			AmountInCents:  1000,
			Currency:       "USD",
		})
		request.Header.Set("Authorization", "Bearer test-token")
		response := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(response, request)

		// Assert
		assert.Equal(t, response.Code, http.StatusBadRequest)
		assert.EqualJSON(t, response.Body.String(), map[string]any{"error": utils.ErrAmountInCentsMismatch.Error()})
	})

	t.Run("invalid request - missing slots", func(t *testing.T) {
		// Arrange
		d1 := insertEmptyPro(t, db)
		defer d1()

		request := NewJsonRequestBody("POST", "/professional/1/booking", booking.CreateBookingRequest{
			IdempotencyKey: "test-key",
			AmountInCents:  1000,
			Currency:       "USD",
		})
		request.Header.Set("Authorization", "Bearer test-token")
		response := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(response, request)

		// Assert
		assert.Equal(t, response.Code, http.StatusBadRequest)
	})

	t.Run("invalid request - invalid currency length", func(t *testing.T) {
		// Arrange
		d1 := insertEmptyPro(t, db)
		defer d1()

		request := NewJsonRequestBody("POST", "/professional/1/booking", booking.CreateBookingRequest{
			Slots: []booking.Slot{{
				Date: "2023-01-01",
				From: "10:00:00",
				To:   "11:00:00",
			}},
			IdempotencyKey: "test-key",
			AmountInCents:  1000,
			Currency:       "US", // Invalid length
		})
		request.Header.Set("Authorization", "Bearer test-token")
		response := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(response, request)

		// Assert
		assert.Equal(t, response.Code, http.StatusBadRequest)
	})

	t.Run("invalid request - missing idempotency key", func(t *testing.T) {
		// Arrange
		d1 := insertEmptyPro(t, db)
		defer d1()

		request := NewJsonRequestBody("POST", "/professional/1/booking", booking.CreateBookingRequest{
			Slots: []booking.Slot{{
				Date: "2023-01-01",
				From: "10:00:00",
				To:   "11:00:00",
			}},
			AmountInCents: 1000,
			Currency:      "USD",
		})
		request.Header.Set("Authorization", "Bearer test-token")
		response := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(response, request)

		// Assert
		assert.Equal(t, response.Code, http.StatusBadRequest)
	})
}

func NewTestBookingController(db *pgxpool.Pool, timeProvider utils.TimeProvider) *booking.BookingController {
	repository := booking.NewRepository(db, timeProvider)
	service := booking.NewService(repository)
	userService := &MockUserService{UserId: 1}
	return booking.NewController(userService, service)
}
