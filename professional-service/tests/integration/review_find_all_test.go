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

func ReviewFindAll(t *testing.T, db *sql.DB) {
	handler := router.Handler(NewTestController(db))

	t.Run("empty review", func(t *testing.T) {
		// Arrange
		professionalId := int64(1)

		d1 := insertPro(t, db, professional.Professional{ID: professionalId})
		defer d1()

		request := NewJsonRequest("GET", fmt.Sprintf("/professional/%d/review?page=%d&pageSize=%d", professionalId, 1, 1), nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Assert
		response_model := []professional.Review{}
		Unmarshal(response, &response_model)
		assert.Equal(t, len(response_model), 0)
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
		professionalId := int64(2)
		date := time.Date(2024, 1, 1, 10, 30, 20, 0, time.UTC)
		d0 := insertUser(t, db, user)
		defer d0()
		d1 := insertPro(t, db, professional.Professional{ID: professionalId})
		defer d1()
		reviews := []professional.Review{
			{
				ID:             67,
				UserID:         int64(user.ID),
				ProfessionalID: professionalId,
				Rate:           4,
				ContentText:    String("It was a good review!"),
				CreatedAt:      date,
				UpdatedAt:      date,
			},
		}
		d2 := insertReview(t, db, reviews...)
		defer d2()

		request := NewJsonRequest("GET", fmt.Sprintf("/professional/%d/review?page=%d&pageSize=%d", professionalId, 1, 1), nil)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)

		// Assert
		assert.Equal(t, response.Code, http.StatusOK)

		response_model := []professional.Review{}
		Unmarshal(response, &response_model)
		assert.Equal(t, len(response_model), 1)

		res := response_model[0]
		assert.Equal(t, res.ID, int64(67))
		assert.Equal(t, res.Rate, 4)
		assert.Equal(t, *res.ContentText, "It was a good review!")
		assert.Equal(t, res.User.Email, user.Email)
		assert.Equal(t, res.User.FirstName, user.FirstName)
		assert.Equal(t, res.User.LastName, user.LastName)
		assert.Equal(t, res.User.ProfileImage, user.ProfileImage)
		assert.Equal(t, res.CreatedAt, date)
		assert.Equal(t, res.UpdatedAt, date)
	})
}
