package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking"
	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/mocks"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BookingCreateTest(t *testing.T, db *pgxpool.Pool) {
	timeProvider := &FakeTimeProvider{}
	userId := 1
	handler := router.Handler(NewTestController(db), NewTestBookingController(db, timeProvider, userId))

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
		d2 := insertUserWithId(t, db, userId)
		defer d2()
		availibility := []professional.Availability{
			{
				ProfessionalID: 1,
				Date:           civil.Date{Year: 2025, Month: 1, Day: 1},
				From:           civil.Time{Hour: 5, Minute: 30},
				To:             civil.Time{Hour: 6, Minute: 30},
			},
		}
		ids, d3 := insertAvailability(t, db, availibility...)
		defer d3()

		request := NewJsonRequestBody("POST", "/professional/1/booking", booking_model.CreateBookingRequest{
			Slots: []booking_model.Slot{
				{Id: ids[0]},
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

		request := NewJsonRequestBody("POST", "/professional/1/booking", booking_model.CreateBookingRequest{
			Slots: []booking_model.Slot{
				{Id: 0},
				{Id: 1},
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

		request := NewJsonRequestBody("POST", "/professional/1/booking", booking_model.CreateBookingRequest{
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

		request := NewJsonRequestBody("POST", "/professional/1/booking", booking_model.CreateBookingRequest{
			Slots: []booking_model.Slot{{
				Id: 1,
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

		request := NewJsonRequestBody("POST", "/professional/1/booking", booking_model.CreateBookingRequest{
			Slots: []booking_model.Slot{{
				Id: 0,
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

func NewTestBookingController(db *pgxpool.Pool, timeProvider utils.TimeProvider, userId int) *booking.BookingController {
	repository := booking.NewRepository(db, timeProvider)
	service := booking.NewService(repository, &mocks.FakePaymentService{})
	userService := &MockUserService{UserId: int64(userId)}
	return booking.NewController(userService, service)
}
