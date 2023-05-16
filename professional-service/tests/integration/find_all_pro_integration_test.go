package integration_test

import (
	"database/sql"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ListProfessionalTest(t *testing.T, db *sql.DB) {
	t.Run("Empty professionals", func(t *testing.T) {
		// Arrange
		request, _ := http.NewRequest("GET", "/professionals", nil)
		response := httptest.NewRecorder()
		sut := router.Handler(db)
		// Act
		sut.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)
		assert.EqualJSON(t, response.Body.String(), []string{})
	})
}
