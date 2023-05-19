package professional

import "context"

type FakeRepository struct {
	findAllSuccess  []Professional
	findAllError    error
	findByIdSuccess Professional
	findByIdError   error
	updateError     error
}

func (r *FakeRepository) FindAll(ctx context.Context, fields ...string) ([]Professional, error) {
	return r.findAllSuccess, r.findAllError
}

func (r *FakeRepository) FindById(ctx context.Context, id string, fields ...string) (Professional, error) {
	return r.findByIdSuccess, r.findByIdError
}

func (r *FakeRepository) Update(ctx context.Context, id string, p UpdateRequest) error {
	return r.updateError
}

type MockUserServiceAlwaysAuthenticated struct {
	IsAuthenticatedCalled bool
}

func (m *MockUserServiceAlwaysAuthenticated) IsAuthenticated(context.Context, string) bool {
	m.IsAuthenticatedCalled = true
	return true
}
