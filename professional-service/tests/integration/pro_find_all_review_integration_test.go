package integration_test

import (
	"database/sql"
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
			ProfileImage: "image.someurl.com",
		}
		proId := int64(2)
		date := time.Date(2024, 1, 1, 10, 30, 20, 0, time.UTC)
		d0 := insertUser(t, db, user)
		defer d0()
		d1 := insertPro(t, db, professional.Professional{ID: proId})
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
}
