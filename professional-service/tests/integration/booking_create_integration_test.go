package integration_test

import (
	"testing"

	mocks2 "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/data/bookingrepo"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/mocks"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BookingCreateTest(t *testing.T, db *pgxpool.Pool) {
	// timeProvider := &FakeTimeProvider{}
	// userId := 1
	// handler := router.Handler(NewTestController(db), NewTestBookingController(db, timeProvider, userId))
}

func NewTestBookingController(db *pgxpool.Pool, timeProvider utils.TimeProvider, userId int) *booking.Controller {
	repository := mocks2.NewRepository(db, timeProvider)
	service := booking.NewService(repository, &mocks.FakePaymentService{})
	userService := &MockUserService{UserId: int64(userId)}
	return booking.NewController(userService, service)
}
