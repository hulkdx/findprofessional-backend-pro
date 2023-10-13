package integration_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"gorm.io/gorm"
)

func FindAllRatingProfessionalTest(t *testing.T, db *sql.DB, gdb *gorm.DB) {
	handler := router.Handler(NewTestController(db))

	t.Run("proffesional got give 5 star ratings from all users", func(t *testing.T) {
		// Arrange
		expected_rating := "5.00"
		now := time.Now()
		professional_records := []professional.Professional{
			{
				ID:        1,
				Email:     "test1@gmail.com",
				Password:  "some_hex_value2",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
		}
		gdb.Create(professional_records)
		defer gdb.Delete(professional_records)

		rating_records := []ProfessionalRating{
			{
				ID:             1,
				UserID:         2,
				ProfessionalID: 1,
				Rate:           5,
			},
			{
				ID:             2,
				UserID:         3,
				ProfessionalID: 1,
				Rate:           5,
			},
			{
				ID:             3,
				UserID:         4,
				ProfessionalID: 1,
				Rate:           5,
			},
			{
				ID:             4,
				UserID:         5,
				ProfessionalID: 1,
				Rate:           5,
			},
		}
		gdb.Create(rating_records)
		defer gdb.Delete(rating_records)

		request := NewJsonRequest("GET", "/professional", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Assert
		assert.Equal(t, response.Code, http.StatusOK)

		rating_response := []RatingResponse{}
		json.Unmarshal(response.Body.Bytes(), &rating_response)
		actual_rating := rating_response[0].Rating
		assert.Equal(t, actual_rating, expected_rating)
	})

	t.Run("proffesional got no ratings", func(t *testing.T) {
		// Arrange
		var expected_rating string

		now := time.Now()
		professional_records := []professional.Professional{
			{
				ID:        1,
				Email:     "test1@gmail.com",
				Password:  "some_hex_value2",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
		}
		gdb.Create(professional_records)
		defer gdb.Delete(professional_records)

		request := NewJsonRequest("GET", "/professional", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Assert
		assert.Equal(t, response.Code, http.StatusOK)

		rating_response := []RatingResponse{}
		json.Unmarshal(response.Body.Bytes(), &rating_response)
		actual_rating := rating_response[0].Rating
		assert.Equal(t, actual_rating, expected_rating)
	})
}

type RatingResponse struct {
	Rating string `json:"rating"`
}
