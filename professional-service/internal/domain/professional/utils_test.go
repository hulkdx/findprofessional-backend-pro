package professional

import "context"

type FakeRepository struct {
	findAllSuccess  []Professional
	findAllError    error
	findByIdSuccess Professional
	findByIdError   error
	createError     error
	deleteError     error
	updateError     error
}

func (r *FakeRepository) FindAll(ctx context.Context, fields ...string) ([]Professional, error) {
	return r.findAllSuccess, r.findAllError
}

func (r *FakeRepository) FindById(ctx context.Context, id string, fields ...string) (Professional, error) {
	return r.findByIdSuccess, r.findByIdError
}

func (r *FakeRepository) Create(ctx context.Context, p Professional) error {
	return r.createError
}
func (r *FakeRepository) Delete(ctx context.Context, id string) error {
	return r.deleteError
}
func (r *FakeRepository) Update(ctx context.Context, id string, p Professional) error {
	return r.updateError
}

type MockUserServiceAlwaysAuthenticated struct{}

func (m *MockUserServiceAlwaysAuthenticated) IsAuthenticated(context.Context, string) bool {
	return true
}
