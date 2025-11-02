package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"github.com/jackc/pgx/v5/pgxpool"
)

func FindAllRatingProfessionalTest(t *testing.T, db *pgxpool.Pool) {
	handler := router.Handler(NewTestController(db), nil)

	t.Run("proffesional got give 5 star ratings from all users", func(t *testing.T) {
		// Arrange
		expected_rating := "5.00"
		professional_records := []model_professional.Professional{
			{
				ID:            1,
				Email:         "test1@gmail.com",
				Password:      "some_hex_value2",
				PriceNumber:   Int(0),
				PriceCurrency: String(""),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}
		d1 := insertPro(t, db, professional_records...)
		defer d1()

		userId := []int{
			2,
			3,
			4,
			5,
		}
		d2 := insertUserWithId(t, db, userId...)
		defer d2()

		rating_records := []model_professional.Review{
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
		d3 := insertReview(t, db, rating_records...)
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

		professional_records := []model_professional.Professional{
			{
				ID:            1,
				Email:         "test1@gmail.com",
				Password:      "some_hex_value2",
				PriceNumber:   Int(0),
				PriceCurrency: String(""),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}
		d1 := insertPro(t, db, professional_records...)
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
