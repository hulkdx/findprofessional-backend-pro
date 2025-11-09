package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks2 "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/data/bookingrepo"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
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
	handler := router.Handler(NewTestController(db), NewTestBookingController(db, timeProvider, userId))

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

func NewTestBookingController(db *pgxpool.Pool, timeProvider utils.TimeProvider, userId int) *booking.Controller {
	repository := mocks2.NewRepository(db, timeProvider)
	service := booking.NewService(repository, &mocks.FakePaymentService{}, timeProvider)
	userService := &MockUserService{UserId: int64(userId)}
	return booking.NewController(userService, service)
}
