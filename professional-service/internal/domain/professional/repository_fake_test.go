package professional

import (
	"context"
)

type FakeRepository struct {
	FindAllSuccess       []Professional
	FindAllError         error
	FindByIdSuccess      Professional
	FindByIdError        error
	UpdateError          error
	FindAllReviewSuccess Reviews
	FindAllReviewError   error
	CreateRequestCalled  CreateRequest
	CreateError          error
	GetBookingStatusResp StatusResponse
	GetBookingStatusErr  error
}

func (r *FakeRepository) FindAll(ctx context.Context) ([]Professional, error) {
	return r.FindAllSuccess, r.FindAllError
}

func (r *FakeRepository) FindById(ctx context.Context, id string) (Professional, error) {
	return r.FindByIdSuccess, r.FindByIdError
}

func (r *FakeRepository) Update(ctx context.Context, id string, p UpdateRequest) error {
	return r.UpdateError
}

func (r *FakeRepository) FindAllReview(ctx context.Context, professionalID int64, page int, pageSize int) (Reviews, error) {
	return r.FindAllReviewSuccess, r.FindAllReviewError
}

func (r *FakeRepository) Create(ctx context.Context, request CreateRequest, pending bool) error {
	r.CreateRequestCalled = request
	return r.CreateError
}

func (m *FakeRepository) GetAvailability(ctx context.Context, professionalId int64) (Availabilities, error) {
	return nil, nil
}

func (m *FakeRepository) UpdateAvailability(ctx context.Context, professionalId int64, availability UpdateAvailabilityRequest) error {
	return nil
}

func (r *FakeRepository) GetBookingStatus(ctx context.Context, bookingId int64) (StatusResponse, error) {
	return r.GetBookingStatusResp, r.GetBookingStatusErr
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
