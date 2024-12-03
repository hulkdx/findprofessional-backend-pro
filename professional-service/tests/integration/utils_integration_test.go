package integration_test

import (
	"context"
	"database/sql"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

func NewTestController(db *sql.DB) *professional.Controller {
	fakeTime := &FakeTimeProvider{
		time.Now(),
	}
	return NewTestControllerWithTimeProvider(db, fakeTime)
}

func NewTestControllerWithTimeProvider(db *sql.DB, timeProvider utils.TimeProvider) *professional.Controller {
	controller := professional.NewController(
		professional.NewService(professional.NewRepository(db, timeProvider)),
		&MockUserService{},
		timeProvider,
	)
	return controller
}

type MockUserService struct{}

func (m *MockUserService) IsAuthenticated(ctx context.Context, auth string) bool {
	return true
}

func (m *MockUserService) Login(ctx context.Context, email string, password string) (string, error) {
	return "", nil
}

func (m *MockUserService) GetAuthenticatedUserId(ctx context.Context, auth string) (*int64, error) {
	return nil, nil
}

type FakeTimeProvider struct {
	NowTime time.Time
}

func (p *FakeTimeProvider) Now() time.Time {
	return p.NowTime
}
