package professional

import "context"

type FakeRepository struct {
	findAllSuccess  []Professional
	findAllError    error
	findByIdSuccess Professional
	findByIdError   error
}

func (r *FakeRepository) FindAll(ctx context.Context, fields ...string) ([]Professional, error) {
	return r.findAllSuccess, r.findAllError
}

func (r *FakeRepository) FindById(ctx context.Context, id string, fields ...string) (Professional, error) {
	return r.findByIdSuccess, r.findByIdError
}

type MockUserServiceAlwaysAuthenticated struct{}

func (m *MockUserServiceAlwaysAuthenticated) IsAuthenticated(context.Context, string) bool {
	return true
}
