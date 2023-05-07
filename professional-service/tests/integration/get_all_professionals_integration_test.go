package integration_test

import (
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListProfessional(t *testing.T) {
	t.Run("Empty professionals", func(t *testing.T) {
		// Arrange
		request, _ := http.NewRequest("GET", "/professional", nil)
		response := httptest.NewRecorder()
		sut := router.Handler()
		// Act
		sut.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, response.Body.String(), "[]")
	})
}