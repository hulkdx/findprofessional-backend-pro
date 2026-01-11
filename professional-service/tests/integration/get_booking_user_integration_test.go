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

func BookingUserTest(t *testing.T, db *pgxpool.Pool) {
	userId := int64(10)
	userService := &MockUserService{UserId: userId}
	handler := router.Handler(NewTestControllerWithUserService(db, userService))

	t.Run("returns bookings for the authenticated user", func(t *testing.T) {
		// Arrange
		otherUserId := int64(11)
		proId := int64(20)
		d1 := insertUser(t, db,
			user.User{
				ID:        userId,
				Email:     "user@email.com",
				FirstName: "User",
				LastName:  "One",
			},
			user.User{
				ID:        otherUserId,
				Email:     "other-user@email.com",
				FirstName: "Other",
				LastName:  "User",
			},
		)
		defer d1()
		d2 := insertPro(t, db, professional.Professional{
			ID:            proId,
			Email:         "pro@email.com",
			FirstName:     "Pro",
			LastName:      "Smith",
			PriceNumber:   Int(30),
			PriceCurrency: String("EUR"),
			Pending:       false,
		})
		defer d2()
		scheduledStart := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		scheduledEnd := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
		bookingID, d3 := insertBooking(t, db, userId, proId, "confirmed", "EUR", "intent-1", &scheduledStart, &scheduledEnd)
		defer d3()
		_, d4 := insertBooking(t, db, otherUserId, proId, "pending", "GBP", "intent-2", &scheduledStart, &scheduledEnd)
		defer d4()
		request := NewJsonRequest("GET", "/professional/bookings/user", nil)
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
		assert.Equal(t, userId, result[0].User.ID)
		assert.Equal(t, proId, result[0].Professional.ID)
	})
}
