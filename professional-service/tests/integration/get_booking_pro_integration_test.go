package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/user"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BookingProTest(t *testing.T, db *pgxpool.Pool) {
	proId := int64(20)
	userService := &MockUserService{UserId: proId}
	handler := router.Handler(NewTestControllerWithUserService(db, userService))

	t.Run("returns bookings for the authenticated professional", func(t *testing.T) {
		// Arrange
		userId := int64(10)
		otherProId := int64(21)
		d1 := insertUser(t, db,
			user.User{
				ID:        userId,
				Email:     "user@email.com",
				FirstName: "User",
				LastName:  "One",
			},
		)
		defer d1()
		d2 := insertPro(t, db,
			professional.Professional{
				ID:            proId,
				Email:         "pro@email.com",
				FirstName:     "Pro",
				LastName:      "Smith",
				PriceNumber:   Int(30),
				PriceCurrency: String("EUR"),
				Pending:       false,
			},
			professional.Professional{
				ID:            otherProId,
				Email:         "other-pro@email.com",
				FirstName:     "Other",
				LastName:      "Pro",
				PriceNumber:   Int(40),
				PriceCurrency: String("GBP"),
				Pending:       false,
			},
		)
		defer d2()
		scheduledStart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		scheduledEnd := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
		bookingID, d3 := insertBooking(t, db, userId, proId, "confirmed", "EUR", "intent-1", &scheduledStart, &scheduledEnd)
		defer d3()
		_, d4 := insertBooking(t, db, userId, otherProId, "pending", "GBP", "intent-2", &scheduledStart, &scheduledEnd)
		defer d4()
		request := NewJsonRequest("GET", "/professional/bookings/pro", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Assert
		assert.Equal(t, http.StatusOK, response.Code)
		var result professional.Bookings
		Unmarshal(response, &result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, bookingID, result[0].ID)
		assert.Equal(t, "confirmed", result[0].Status)
		assert.Equal(t, "EUR", result[0].Currency)
		assert.Equal(t, proId, result[0].Professional.ID)
		assert.Equal(t, userId, result[0].User.ID)
	})

	t.Run("should return session informations from professionals when booking.session_info is empty", func(t *testing.T) {
		// Arrange
		userId := int64(10)
		sessionLink := "https://meet.google.com/abc-defg-hij"
		sessionPlatform := "Google Meet"

		d1 := insertUser(t, db,
			user.User{
				ID:        userId,
				Email:     "user@email.com",
				FirstName: "User",
				LastName:  "One",
			},
		)
		defer d1()
		d2 := insertPro(t, db,
			professional.Professional{
				ID:              proId,
				Email:           "pro@email.com",
				FirstName:       "Pro",
				LastName:        "Smith",
				PriceNumber:     Int(30),
				PriceCurrency:   String("EUR"),
				Pending:         false,
				SessionLink:     &sessionLink,
				SessionPlatform: &sessionPlatform,
			},
		)
		defer d2()
		scheduledStart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		scheduledEnd := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
		_, d3 := insertBooking(t, db, userId, proId, "confirmed", "EUR", "intent-1", &scheduledStart, &scheduledEnd)
		defer d3()
		request := NewJsonRequest("GET", "/professional/bookings/pro", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Assert
		assert.Equal(t, http.StatusOK, response.Code)
		var result professional.Bookings
		Unmarshal(response, &result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, *result[0].Session.Link, sessionLink)
		assert.Equal(t, *result[0].Session.Platform, sessionPlatform)
	})

	t.Run("should return session informations from bookings when booking.session_info is not empty", func(t *testing.T) {
		// Arrange
		userId := int64(10)
		proSessionLink := "https://meet.google.com/abc-defg-hij"
		proSessionPlatform := "Google Meet"

		bookingSessionLink := "https://zoom.google.com/abc-defg-hij"
		bookingSessionPlatform := "Zoom"

		d1 := insertUser(t, db,
			user.User{
				ID:        userId,
				Email:     "user@email.com",
				FirstName: "User",
				LastName:  "One",
			},
		)
		defer d1()
		d2 := insertPro(t, db,
			professional.Professional{
				ID:              proId,
				Email:           "pro@email.com",
				FirstName:       "Pro",
				LastName:        "Smith",
				PriceNumber:     Int(30),
				PriceCurrency:   String("EUR"),
				Pending:         false,
				SessionLink:     &proSessionLink,
				SessionPlatform: &proSessionPlatform,
			},
		)
		defer d2()
		scheduledStart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		scheduledEnd := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
		_, d3 := insertBookingWithSessions(
			t,
			db,
			userId,
			proId,
			"confirmed",
			"EUR",
			"intent-1",
			&scheduledStart,
			&scheduledEnd,
			&bookingSessionLink,
			&bookingSessionPlatform,
		)
		defer d3()
		request := NewJsonRequest("GET", "/professional/bookings/pro", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Assert
		assert.Equal(t, http.StatusOK, response.Code)
		var result professional.Bookings
		Unmarshal(response, &result)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, *result[0].Session.Link, bookingSessionLink)
		assert.Equal(t, *result[0].Session.Platform, bookingSessionPlatform)
	})
}
