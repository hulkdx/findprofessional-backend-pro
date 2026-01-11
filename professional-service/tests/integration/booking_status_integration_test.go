package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/user"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BookingStatusTest(t *testing.T, db *pgxpool.Pool) {
	userId := int64(10)
	userService := &MockUserService{UserId: userId}
	handler := router.Handler(NewTestControllerWithUserService(db, userService))

	t.Run("booking not found", func(t *testing.T) {
		// Arrange
		request := NewJsonRequest("GET", "/professional/booking/999/status", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Assert
		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("booking found", func(t *testing.T) {
		// Arrange
		proId := int64(20)

		d1 := insertUser(t, db, user.User{
			ID:    userId,
			Email: "user@email.com",
		})
		defer d1()
		d2 := insertPro(t, db, professional.Professional{
			ID:            proId,
			Email:         "pro@email.com",
			PriceNumber:   Int(30),
			PriceCurrency: String("EUR"),
			Pending:       false,
		})
		defer d2()
		bookingID, d3 := insertBooking(
			t,
			db,
			userId,
			proId,
			"confirmed",
			"EUR",
			"intent-1",
			nil,
			nil,
		)
		defer d3()
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, NewBookingStatusRequest(bookingID))
		// Assert
		assert.Equal(t, http.StatusOK, response.Code)
		var result professional.StatusResponse
		Unmarshal(response, &result)
		assert.Equal(t, "confirmed", result.Status)
	})
	t.Run("user should only see their own bookings", func(t *testing.T) {
		// Arrange
		otherUserId := int64(11)
		proId := int64(20)
		d1 := insertUser(t, db,
			user.User{
				ID:    userId,
				Email: "user@email.com",
			},
			user.User{
				ID:    otherUserId,
				Email: "other-user@email.com",
			},
		)
		defer d1()
		d2 := insertPro(t, db,
			professional.Professional{
				ID:            proId,
				Email:         "pro@email.com",
				PriceNumber:   Int(0),
				PriceCurrency: String("GBP"),
				Pending:       false,
			},
		)
		defer d2()
		bookingID, d3 := insertBooking(t, db, otherUserId, proId, "confirmed", "EUR", "intent-1", nil, nil)
		defer d3()
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, NewBookingStatusRequest(bookingID))
		// Assert
		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func NewBookingStatusRequest(bookingID int64) *http.Request {
	return NewJsonRequest("GET", fmt.Sprintf("/professional/booking/%d/status", bookingID), nil)
}
