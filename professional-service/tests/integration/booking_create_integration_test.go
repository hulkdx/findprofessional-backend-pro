package integration_test

import (
	"fmt"
	"net/http"
	"testing"

	mocks2 "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/data/bookingrepo"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/mocks"
	"github.com/jackc/pgx/v5/pgxpool"
)

// noinspection GoUnhandledErrorResult
func BookingCreateTest(t *testing.T, db *pgxpool.Pool) {
	/* TODO:
	timeProvider := &mocks.FakeTimeProvider{}
	userId := 1
	handler := router.Handler(NewTestController(db), NewTestBookingController(db, timeProvider, userId))

	t.Run("idempotency key conflict returns conflict status", func(t *testing.T) {
		// Arrange
		idempotencyKey := "019a4eb2-edb4-7068-b1cb-89bbaee6581f"
		proId := 1
		requestBody := bookingmodel.CreateBookingRequest{
			Availabilities: []bookingmodel.Availability{{Id: 1}},
			IdempotencyKey: idempotencyKey,
			AmountInCents:  1500,
			Currency:       "GBP",
		}

		d1 := insertUserWithId(t, db, userId)
		defer d1()
		d2 := insertBookingHolds(t, db, userId, idempotencyKey, timeProvider.Now().UTC(), timeProvider.Now().UTC().Add(1*time.Hour))
		defer d2()

		request := CreateBookingRequest(proId, requestBody)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Assert
		assert.Equal(t, response.Code, http.StatusConflict)
	})
	*/
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
