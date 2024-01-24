package professional

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func TestFindAllReview(t *testing.T) {
	t.Run("authorize", func(t *testing.T) {
		// Arrange
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}
		// Act
		controller.FindAllReview(httptest.NewRecorder(), findAllReviewRequest("1", "1", "1"))
		// Assert
		assert.Equal(t, userService.IsAuthenticatedCalled, true)
	})
}

func TestFindAllReviewParameters(t *testing.T) {
	t.Run("id is int then no invalid id error", func(t *testing.T) {
		// Arrange
		id := "1"
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}

		response := httptest.NewRecorder()
		// Act
		controller.FindAllReview(response, findAllReviewRequest(id, "1", "1"))
		// Assert
		assertNotError(t, response, "id is in wrong format")
	})
	t.Run("id is not int then error", func(t *testing.T) {
		// Arrange
		id := "a"
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}

		response := httptest.NewRecorder()
		// Act
		controller.FindAllReview(response, findAllReviewRequest(id, "1", "1"))
		// Assert
		assertError(t, response, "id is in wrong format")
	})
	t.Run("page is int then no invalid id error", func(t *testing.T) {
		// Arrange
		page := "1"
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}

		response := httptest.NewRecorder()
		// Act
		controller.FindAllReview(response, findAllReviewRequest("", page, "1"))
		// Assert
		assertNotError(t, response, "page is in wrong format")
	})
	t.Run("page is not int then error", func(t *testing.T) {
		// Arrange
		page := "a"
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}

		response := httptest.NewRecorder()
		// Act
		controller.FindAllReview(response, findAllReviewRequest("1", page, "1"))
		// Assert
		assertError(t, response, "page is in wrong format")
	})
	t.Run("pageSize is int then no invalid id error", func(t *testing.T) {
		// Arrange
		pageSize := "1"
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}

		response := httptest.NewRecorder()
		// Act
		controller.FindAllReview(response, findAllReviewRequest("1", "1", pageSize))
		// Assert
		assertNotError(t, response, "pageSize is in wrong format")
	})
	t.Run("pageSize is not int then error", func(t *testing.T) {
		// Arrange
		pageSize := "a"
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}

		response := httptest.NewRecorder()
		// Act
		controller.FindAllReview(response, findAllReviewRequest("1", "1", pageSize))
		// Assert
		assertError(t, response, "pageSize is in wrong format")
	})
	t.Run("page cannot be zero", func(t *testing.T) {
		// Arrange
		page := "0"
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}

		response := httptest.NewRecorder()
		// Act
		controller.FindAllReview(response, findAllReviewRequest("1", page, "1"))
		// Assert
		assertError(t, response, "page is in wrong format")
	})
	t.Run("page cannot be negative", func(t *testing.T) {
		// Arrange
		page := "-1"
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}

		response := httptest.NewRecorder()
		// Act
		controller.FindAllReview(response, findAllReviewRequest("1", page, "1"))
		// Assert
		assertError(t, response, "page is in wrong format")
	})
	t.Run("pageSize cannot be zero", func(t *testing.T) {
		// Arrange
		pageSize := "0"
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}

		response := httptest.NewRecorder()
		// Act
		controller.FindAllReview(response, findAllReviewRequest("1", "1", pageSize))
		// Assert
		assertError(t, response, "pageSize is in wrong format")
	})
	t.Run("pageSize cannot be negative", func(t *testing.T) {
		// Arrange
		pageSize := "-1"
		userService := &MockUserServiceAlwaysAuthenticated{}
		controller := &Controller{
			service:     NewService(&FakeRepository{}),
			userService: userService,
		}

		response := httptest.NewRecorder()
		// Act
		controller.FindAllReview(response, findAllReviewRequest("1", "1", pageSize))
		// Assert
		assertError(t, response, "pageSize is in wrong format")
	})
}

func assertError(t *testing.T, response *httptest.ResponseRecorder, err string) {
	assert.EqualJSON(t, response.Body.String(), map[string]any{
		"error": err,
	})
}

func assertNotError(t *testing.T, response *httptest.ResponseRecorder, err string) {
	assert.NotEqualJSON(t, response.Body.String(), map[string]any{
		"error": err,
	})
}

func findAllReviewRequest(id string, page string, pageSize string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)

	url := fmt.Sprintf("/professionals/%s/review?page=%s&pageSize=%s", id, page, pageSize)
	request, _ := http.NewRequest("GET", url, nil)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
	return request
}
