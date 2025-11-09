package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/data/bookingrepo"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/payment"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/mocks"
	"github.com/jackc/pgx/v5/pgxpool"
)

// noinspection GoUnhandledErrorResult
func BookingCreateTest(t *testing.T, db *pgxpool.Pool) {
	timeProvider := &mocks.FakeTimeProvider{}
	userId := 1
	paymentService := mocks.FakePaymentService{}
	handler := router.Handler(NewTestController(db), NewTestBookingController(db, timeProvider, userId, &paymentService))

	t.Run("success", func(t *testing.T) {
		// Arrange
		proId := 1
		idempotencyKey := "019a4eb2-edb4-7068-b1cb-89bbaee6581f"

		d1 := insertUserWithId(t, db, userId)
		defer d1()

		d2 := insertPro(t, db, professional.Professional{ID: int64(proId)})
		defer d2()

		availability := professional.Availability{
			ProfessionalID: int64(proId),
			From:           time.Now().Add(time.Hour),
			To:             time.Now().Add(2 * time.Hour),
		}
		availabilityIds, d3 := insertAvailability(t, db, availability)
		defer d3()

		requestBody := bookingmodel.CreateBookingRequest{
			Availabilities: []bookingmodel.Availability{{Id: availabilityIds[0]}},
			IdempotencyKey: idempotencyKey,
			AmountInCents:  1500,
			Currency:       "EUR",
		}

		request := CreateBookingRequest(proId, requestBody)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Assert
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, response.Body.String(), http.StatusOK)
	})

	t.Run("when availability_id doesn't exists return error 404", func(t *testing.T) {
		// Arrange
		idempotencyKey := "019a4eb2-edb4-7068-b1cb-89bbaee6581f"
		proId := 1
		availabilityId := int64(1)

		requestBody := bookingmodel.CreateBookingRequest{
			Availabilities: []bookingmodel.Availability{{Id: availabilityId}},
			IdempotencyKey: idempotencyKey,
			AmountInCents:  1500,
			Currency:       "EUR",
		}

		d1 := insertUserWithId(t, db, userId)
		defer d1()

		request := CreateBookingRequest(proId, requestBody)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Assert
		assert.Equal(t, response.Code, http.StatusNotFound)
		assert.EqualJSON(t, response.Body.String(), map[string]any{
			"error": utils.ErrAvailabilityDoesNotExist.Error(),
		})
	})
}

func CreateBookingRequest(proId int, requestBody bookingmodel.CreateBookingRequest) *http.Request {
	path := fmt.Sprintf("/professional/%d/booking", proId)
	request := NewJsonRequestBody("POST", path, requestBody)
	return request
}

func NewTestBookingController(
	db *pgxpool.Pool,
	timeProvider utils.TimeProvider,
	userId int,
	paymentService payment.Service,
) *booking.Controller {
	userService := &MockUserService{UserId: int64(userId)}

	repository := bookingrepo.NewRepository(db, timeProvider)
	service := booking.NewService(repository, paymentService, timeProvider)
	return booking.NewController(userService, service)
}
