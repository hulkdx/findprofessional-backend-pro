package professional

import (
	"context"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func TestCreate(t *testing.T) {
	t.Run("should bcrypt password", func(t *testing.T) {
		// Arrange
		rawPassword := "password"
		request := CreateRequest{
			Email:     "test@gmail.com",
			Password:  rawPassword,
			FirstName: "",
			LastName:  "",
			SkypeId:   "",
			AboutMe:   "",
		}
		repository := &FakeRepository{}
		service := NewService(repository)
		// Act
		service.Create(context.Background(), request)
		// Assert
		calledPassword := repository.CreateRequestCalled.Password
		assert.NotEqual(t, calledPassword, rawPassword)
	})
}
