package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"github.com/jackc/pgx/v5/pgxpool"
)

func FindProfessionalTest(t *testing.T, db *pgxpool.Pool) {
	handler := router.Handler(NewTestController(db), nil)

	t.Run("Empty database", func(t *testing.T) {
		// Arrange
		id := 1
		request := NewJsonRequest("GET", fmt.Sprintf("/professional/%d", id), nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	t.Run("found a record", func(t *testing.T) {
		// Arrange
		id := int64(1)
		record := &model_professional.Professional{
			ID:            id,
			Email:         "emailofidone@email.com",
			PriceNumber:   Int(0),
			PriceCurrency: String(""),
		}
		request := NewJsonRequest("GET", fmt.Sprintf("/professional/%d", id), nil)
		response := httptest.NewRecorder()
		d1 := insertPro(t, db, *record)
		defer d1()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)
		assert.EqualJSON(t, response.Body.String(), record)
	})
}
