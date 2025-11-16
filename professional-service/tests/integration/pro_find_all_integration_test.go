package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"github.com/jackc/pgx/v5/pgxpool"
)

func FindAllProfessionalTest(t *testing.T, db *pgxpool.Pool) {
	handler := router.Handler(NewTestController(db))

	t.Run("Empty professionals", func(t *testing.T) {
		// Arrange
		request := NewJsonRequest("GET", "/professional", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)
		assert.EqualJSON(t, response.Body.String(), []string{})
	})

	t.Run("some professional, only show valid data", func(t *testing.T) {
		// Arrange
		records := []professional.Professional{
			{
				ID:            1,
				Email:         "test1@gmail.com",
				Password:      "some_hex_value2",
				PriceNumber:   Int(0),
				PriceCurrency: String(""),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			{
				ID:            2,
				Email:         "test2@gmail.com",
				Password:      "some_hex_value2",
				PriceNumber:   Int(0),
				PriceCurrency: String(""),
			},
		}
		d1 := insertPro(t, db, records...)
		defer d1()
		expected := []professional.Professional{
			{
				ID:            1,
				Email:         "test1@gmail.com",
				PriceNumber:   Int(0),
				PriceCurrency: String(""),
			},
			{
				ID:            2,
				Email:         "test2@gmail.com",
				PriceNumber:   Int(0),
				PriceCurrency: String(""),
			},
		}
		request := NewJsonRequest("GET", "/professional", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Assert
		assert.Equal(t, response.Code, http.StatusOK)
		assert.EqualAnyOrderJSON(t, response.Body.String(), expected)
	})

	t.Run("hide professional with null PriceNumber or PriceCurrency", func(t *testing.T) {
		// Arrange
		records := []professional.Professional{
			{
				ID:            1,
				Email:         "test1@gmail.com",
				Password:      "some_hex_value2",
				PriceNumber:   Int(0),
				PriceCurrency: String("USD"),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			{
				ID:            2,
				Email:         "test2@gmail.com",
				Password:      "some_hex_value2",
				PriceNumber:   nil,
				PriceCurrency: String(""),
			},
			{
				ID:            3,
				Email:         "test3@gmail.com",
				Password:      "some_hex_value3",
				PriceNumber:   Int(50),
				PriceCurrency: nil,
			},
		}
		d1 := insertPro(t, db, records...)
		defer d1()
		expected := []professional.Professional{
			{
				ID:            1,
				Email:         "test1@gmail.com",
				PriceNumber:   Int(0),
				PriceCurrency: String("USD"),
			},
		}
		request := NewJsonRequest("GET", "/professional", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Assert
		assert.Equal(t, response.Code, http.StatusOK)
		assert.EqualAnyOrderJSON(t, response.Body.String(), expected)
	})
}
