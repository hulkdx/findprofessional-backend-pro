package integration_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/user"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func FindAllReviewProfessionalTest(t *testing.T, db *sql.DB) {
	handler := router.Handler(NewTestController(db))

	t.Run("empty review", func(t *testing.T) {
		// Arrange
		d1 := insertEmptyPro(t, db)
		defer d1()

		request := NewJsonRequest("GET", "/professional", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Assert
		assert.Equal(t, response.Code, http.StatusOK)

		response_model := []professional.Professional{}
		Unmarshal(response, &response_model)
		assert.Equal(t, len(response_model), 1)
		assert.Equal(t, len(response_model[0].Review), 0)
	})

	t.Run("some review", func(t *testing.T) {
		// Arrange
		user := user.User{
			ID:           1,
			Email:        "test_user@gmail.com",
			FirstName:    "user first name",
			LastName:     "user last name",
			ProfileImage: String("image.someurl.com"),
		}
		proId := int64(2)
		pro := professional.Professional{
			ID:            proId,
			PriceNumber:   Int(0),
			PriceCurrency: String(""),
		}
		date := time.Date(2024, 1, 1, 10, 30, 20, 0, time.UTC)
		d0 := insertUser(t, db, user)
		defer d0()
		d1 := insertPro(t, db, pro)
		defer d1()
		reviews := []professional.Review{
			{
				ID:             67,
				UserID:         int64(user.ID),
				ProfessionalID: proId,
				Rate:           4,
				ContentText:    String("It was a good review!"),
				CreatedAt:      date,
				UpdatedAt:      date,
			},
		}
		d2 := insertReview(t, db, reviews...)
		defer d2()

		request := NewJsonRequest("GET", "/professional", nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)

		// Assert
		assert.Equal(t, response.Code, http.StatusOK)

		response_model := []professional.Professional{}
		Unmarshal(response, &response_model)
		assert.Equal(t, len(response_model), 1)
		assert.Equal(t, response_model[0].ID, proId)
		assert.Equal(t, len(response_model[0].Review), 1)
		assert.Equal(t, response_model[0].Review[0].ID, int64(67))
		assert.Equal(t, response_model[0].Review[0].Rate, 4)
		assert.Equal(t, *response_model[0].Review[0].ContentText, "It was a good review!")
		assert.Equal(t, response_model[0].Review[0].User.Email, user.Email)
		assert.Equal(t, response_model[0].Review[0].User.FirstName, user.FirstName)
		assert.Equal(t, response_model[0].Review[0].User.LastName, user.LastName)
		assert.Equal(t, response_model[0].Review[0].User.ProfileImage, user.ProfileImage)
		assert.Equal(t, response_model[0].Review[0].CreatedAt, date)
		assert.Equal(t, response_model[0].Review[0].UpdatedAt, date)
	})

	t.Run("limit reviews by 3", func(t *testing.T) {
		testCases := []struct {
			reviewCount        int
			expectedReviews    int
			expectedReviewSize int64
		}{
			{
				reviewCount:        1,
				expectedReviews:    1,
				expectedReviewSize: 1,
			},
			{
				reviewCount:        2,
				expectedReviews:    2,
				expectedReviewSize: 2,
			},
			{
				reviewCount:        3,
				expectedReviews:    3,
				expectedReviewSize: 3,
			},
			{
				reviewCount:        4,
				expectedReviews:    3,
				expectedReviewSize: 4,
			},
			{
				reviewCount:        5,
				expectedReviews:    3,
				expectedReviewSize: 5,
			},
		}
		for _, tc := range testCases {
			t.Run(fmt.Sprintf("when %d reviews in the db, api response has %d review count", tc.reviewCount, tc.expectedReviews), func(t *testing.T) {
				// Arrange
				d := createReviews(t, db, tc.reviewCount)
				defer d()
				request := NewJsonRequest("GET", "/professional", nil)
				response := httptest.NewRecorder()
				// Act
				handler.ServeHTTP(response, request)
				// Assert
				response_model := []professional.Professional{}
				Unmarshal(response, &response_model)
				assert.Equal(t, len(response_model[0].Review), tc.expectedReviews)
				assert.Equal(t, response_model[0].ReviewSize, tc.expectedReviewSize)
			})
		}
	})

}

func createReviews(t *testing.T, db *sql.DB, count int) func() {
	proId := int64(2)
	pro := professional.Professional{
		ID:            proId,
		PriceNumber:   Int(0),
		PriceCurrency: String(""),
	}
	d1 := insertPro(t, db, pro)

	userIds := []int{}
	for i := 1; i <= count; i++ {
		userIds = append(userIds, i)
	}
	d2 := insertUserWithId(t, db, userIds...)

	var reviews []professional.Review
	for i := 0; i < count; i++ {
		review := professional.Review{
			ID:             int64(i),
			UserID:         int64(userIds[i]),
			ProfessionalID: proId,
			Rate:           4,
			ContentText:    String("It was a good review!"),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		reviews = append(reviews, review)
	}
	d3 := insertReview(t, db, reviews...)

	return func() {
		d3()
		d2()
		d1()
	}
}
