package professional

import "context"

type FakeRepository struct {
	findAllSuccess       []Professional
	findAllError         error
	findByIdSuccess      Professional
	findByIdError        error
	updateError          error
	findAllReviewSuccess Reviews
	findAllReviewError   error
	createRequestCalled  CreateRequest
	createError          error
}

func (r *FakeRepository) FindAll(ctx context.Context) ([]Professional, error) {
	return r.findAllSuccess, r.findAllError
}

func (r *FakeRepository) FindById(ctx context.Context, id string) (Professional, error) {
	return r.findByIdSuccess, r.findByIdError
}

func (r *FakeRepository) Update(ctx context.Context, id string, p UpdateRequest) error {
	return r.updateError
}

func (r *FakeRepository) FindAllReview(ctx context.Context, professionalID int64, page int, pageSize int) (Reviews, error) {
	return r.findAllReviewSuccess, r.findAllReviewError
}

func (r *FakeRepository) Create(ctx context.Context, request CreateRequest, pending bool) error {
	r.createRequestCalled = request
	return r.createError
}

type MockUserServiceAlwaysAuthenticated struct {
	IsAuthenticatedCalled bool
}

func (m *MockUserServiceAlwaysAuthenticated) IsAuthenticated(context.Context, string) bool {
	m.IsAuthenticatedCalled = true
	return true
}

func (m *MockUserServiceAlwaysAuthenticated) Login(ctx context.Context, email string, password string) (string, error) {
	return "", nil
}

func (m *MockUserServiceAlwaysAuthenticated) GetAuthenticatedUserId(ctx context.Context, auth string) (*int64, error) {
	return nil, nil
}
