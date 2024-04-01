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
	createError          error
}

func (r *FakeRepository) FindAll(ctx context.Context, filterQuery string, filterItems FilterItems) ([]Professional, error) {
	return r.findAllSuccess, r.findAllError
}

func (r *FakeRepository) FindById(ctx context.Context, id string, filterQuery string, filterItems FilterItems) (Professional, error) {
	return r.findByIdSuccess, r.findByIdError
}

func (r *FakeRepository) Update(ctx context.Context, id string, p UpdateRequest) error {
	return r.updateError
}

func (r *FakeRepository) FindAllReview(ctx context.Context, professionalID int64, page int, pageSize int) (Reviews, error) {
	return r.findAllReviewSuccess, r.findAllReviewError
}

func (r *FakeRepository) Create(ctx context.Context, request CreateRequest) error {
	return r.createError
}

type MockUserServiceAlwaysAuthenticated struct {
	IsAuthenticatedCalled bool
}

func (m *MockUserServiceAlwaysAuthenticated) IsAuthenticated(context.Context, string) bool {
	m.IsAuthenticatedCalled = true
	return true
}
