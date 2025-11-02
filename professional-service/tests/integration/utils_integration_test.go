package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/data/professionalrepo"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/user"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewTestController(db *pgxpool.Pool) *professional.Controller {
	fakeTime := &FakeTimeProvider{
		time.Now(),
	}
	return NewTestControllerWithTimeProvider(db, fakeTime)
}

func NewTestControllerWithUserService(db *pgxpool.Pool, userService user.Service) *professional.Controller {
	fakeTime := &FakeTimeProvider{
		time.Now(),
	}
	return NewTestControllerWithTimeProviderWithUserService(db, fakeTime, userService)
}

func NewTestControllerWithTimeProvider(db *pgxpool.Pool, timeProvider utils.TimeProvider) *professional.Controller {
	return NewTestControllerWithTimeProviderWithUserService(db, timeProvider, &MockUserService{})
}

func NewTestControllerWithTimeProviderWithUserService(db *pgxpool.Pool, timeProvider utils.TimeProvider, userService user.Service) *professional.Controller {
	controller := professional.NewController(
		professional.NewService(professionalrepo.NewRepository(db, timeProvider)),
		userService,
		timeProvider,
	)
	return controller
}

type MockUserService struct {
	UserId int64
}

func (m *MockUserService) IsAuthenticated(ctx context.Context, auth string) bool {
	return true
}

func (m *MockUserService) Login(ctx context.Context, email string, password string) (string, error) {
	return "", nil
}

func (m *MockUserService) GetAuthenticatedUserId(ctx context.Context, auth string) (int64, error) {
	return m.UserId, nil
}

type FakeTimeProvider struct {
	NowTime time.Time
}

func (p *FakeTimeProvider) Now() time.Time {
	return p.NowTime
}

func NewJsonRequest(method, url string, body io.Reader) *http.Request {
	request, _ := http.NewRequest(method, url, body)
	request.Header.Set("Content-Type", "application/json")
	return request
}

func NewJsonRequestBody(method, url string, body any) *http.Request {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	return NewJsonRequest(method, url, &buf)
}
