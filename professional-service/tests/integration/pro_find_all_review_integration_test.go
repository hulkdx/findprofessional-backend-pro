package integration_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func FindAllReviewProfessionalTest(t *testing.T, db *sql.DB) {
	handler := router.Handler(NewTestController(db))

	t.Run("empty review", func(t *testing.T) {
		// Arrange
		d1 := insertEmptyPro(db)
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
		assert.Equal(t, response_model[0].Review.Total, 0)
		assert.Equal(t, len(response_model[0].Review.Content), 0)
	})
}
