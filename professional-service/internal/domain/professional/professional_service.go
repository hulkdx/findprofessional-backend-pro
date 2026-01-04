package professional

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	FindAll(context.Context) ([]Professional, error)
	FindById(ctx context.Context, id string) (Professional, error)
	Create(context.Context, CreateRequest) error
	Update(ctx context.Context, id string, p UpdateRequest) error
	FindAllReview(ctx context.Context, professionalId int64, page int, pageSize int) (Reviews, error)
	GetAvailability(ctx context.Context, professionalId int64) (Availabilities, error)
	UpdateAvailability(ctx context.Context, professionalId int64, availability UpdateAvailabilityRequest) error
	GetBookingStatus(ctx context.Context, bookingId int64, userId int64) (StatusResponse, error)
	GetBookings(ctx context.Context, id int64, idType UserType) (*Bookings, error)
}

type serviceImpl struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &serviceImpl{repository}
}

func (s *serviceImpl) FindAll(ctx context.Context) ([]Professional, error) {
	return s.repository.FindAll(ctx)
}

func (s *serviceImpl) FindById(ctx context.Context, id string) (Professional, error) {
	return s.repository.FindById(ctx, id)
}

func (s *serviceImpl) Create(ctx context.Context, r CreateRequest) error {
	pending := true
	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	r.Password = string(hash)

	err = s.repository.Create(ctx, r, pending)

	var pqErr *pgconn.PgError
	if errors.As(err, &pqErr) && pqErr.Code == pgerrcode.UniqueViolation {
		return utils.ErrDuplicate
	}
	if errors.Is(err, sql.ErrNoRows) {
		return utils.ErrNotFoundUser
	}
	return err
}

func (s *serviceImpl) Update(ctx context.Context, id string, p UpdateRequest) error {
	return s.repository.Update(ctx, id, p)
}

func (s *serviceImpl) FindAllReview(ctx context.Context, professionalId int64, page int, pageSize int) (Reviews, error) {
	return s.repository.FindAllReview(ctx, professionalId, page, pageSize)
}

func (s *serviceImpl) GetAvailability(ctx context.Context, professionalId int64) (Availabilities, error) {
	return s.repository.GetAvailability(ctx, professionalId)
}

func (s *serviceImpl) UpdateAvailability(ctx context.Context, professionalId int64, availability UpdateAvailabilityRequest) error {
	return s.repository.UpdateAvailability(ctx, professionalId, availability)
}

func (s *serviceImpl) GetBookingStatus(ctx context.Context, bookingId int64, userId int64) (StatusResponse, error) {
	return s.repository.GetBookingStatus(ctx, bookingId, userId)
}

func (s *serviceImpl) GetBookings(ctx context.Context, id int64, idType UserType) (*Bookings, error) {
	return s.repository.GetBookings(ctx, id, idType)
}
