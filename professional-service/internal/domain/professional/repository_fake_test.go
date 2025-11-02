package professional

import (
	"context"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional/model"
)

type FakeRepository struct {
	findAllSuccess       []model_professional.Professional
	findAllError         error
	findByIdSuccess      model_professional.Professional
	findByIdError        error
	updateError          error
	findAllReviewSuccess model_professional.Reviews
	findAllReviewError   error
	createRequestCalled  CreateRequest
	createError          error
}

func (r *FakeRepository) FindAll(ctx context.Context) ([]model_professional.Professional, error) {
	return r.findAllSuccess, r.findAllError
}

func (r *FakeRepository) FindById(ctx context.Context, id string) (model_professional.Professional, error) {
	return r.findByIdSuccess, r.findByIdError
}

func (r *FakeRepository) Update(ctx context.Context, id string, p UpdateRequest) error {
	return r.updateError
}

func (r *FakeRepository) FindAllReview(ctx context.Context, professionalID int64, page int, pageSize int) (model_professional.Reviews, error) {
	return r.findAllReviewSuccess, r.findAllReviewError
}

func (r *FakeRepository) Create(ctx context.Context, request CreateRequest, pending bool) error {
	r.createRequestCalled = request
	return r.createError
}

func (m *FakeRepository) GetAvailability(ctx context.Context, professionalId int64) (model_professional.Availabilities, error) {
	return nil, nil
}

func (m *FakeRepository) UpdateAvailability(ctx context.Context, professionalId int64, availability UpdateAvailabilityRequest) error {
	return nil
}

// --------------------------------

type MockUserServiceAlwaysAuthenticated struct {
	IsAuthenticatedCalled        bool
	GetAuthenticatedUserIdCalled bool
}

func (m *MockUserServiceAlwaysAuthenticated) IsAuthenticated(context.Context, string) bool {
	m.IsAuthenticatedCalled = true
	return true
}

func (m *MockUserServiceAlwaysAuthenticated) Login(ctx context.Context, email string, password string) (string, error) {
	return "", nil
}

func (m *MockUserServiceAlwaysAuthenticated) GetAuthenticatedUserId(ctx context.Context, auth string) (int64, error) {
	m.GetAuthenticatedUserIdCalled = true
	return -2, nil
}

func (m *MockUserServiceAlwaysAuthenticated) UpdateAvailability(ctx context.Context, professionalId int64, availability UpdateAvailabilityRequest) error {
	return nil
}
