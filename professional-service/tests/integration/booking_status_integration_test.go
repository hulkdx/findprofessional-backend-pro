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
	handler := router.Handler(NewTestController(db))

	t.Run("booking not found", func(t *testing.T) {
		request := NewJsonRequest("GET", "/professional/booking/999/status", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("booking found", func(t *testing.T) {
		handler := router.Handler(NewTestController(db))
		userCleanup := insertUser(t, db, user.User{
			ID:    10,
			Email: "user@email.com",
		})
		proCleanup := insertPro(t, db, professional.Professional{
			ID:            20,
			Email:         "pro@email.com",
			PriceNumber:   Int(0),
			PriceCurrency: String("GBP"),
			Pending:       false,
		})
		bookingID, bookingCleanup := insertBooking(t, db, 10, 20, "confirmed", "GBP", "intent-1")
		defer bookingCleanup()
		defer proCleanup()
		defer userCleanup()

		request := NewJsonRequest("GET", fmt.Sprintf("/professional/booking/%d/status", bookingID), nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		var result professional.StatusResponse
		Unmarshal(response, &result)
		assert.Equal(t, "confirmed", result.Status)
	})
}
