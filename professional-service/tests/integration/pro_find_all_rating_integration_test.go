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
)

func FindAllRatingProfessionalTest(t *testing.T, db *sql.DB) {
	handler := router.Handler(NewTestController(db))

	t.Run("proffesional got give 5 star ratings from all users", func(t *testing.T) {
		// Arrange
		expected_rating := "5.00"
		professional_records := []professional.Professional{
			{
				ID:        1,
				Email:     "test1@gmail.com",
				Password:  "some_hex_value2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		d1 := insertPro(db, professional_records...)
		defer d1()

		d2 := insertUser(db, []User{
			{ID: 2},
			{ID: 3},
			{ID: 4},
			{ID: 5},
		}...)
		defer d2()

		rating_records := []professional.ProfessionalRating{
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
		d3 := insertRating(db, rating_records...)
		defer d3()

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

		professional_records := []professional.Professional{
			{
				ID:        1,
				Email:     "test1@gmail.com",
				Password:  "some_hex_value2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		d1 := insertPro(db, professional_records...)
		defer d1()

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
