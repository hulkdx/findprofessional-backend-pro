package professional

import (
	"context"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	FindAll(context.Context) ([]Professional, error)
	FindById(ctx context.Context, id string) (Professional, error)
	Create(context.Context, CreateRequest) error
	Update(ctx context.Context, id string, p UpdateRequest) error
	FindAllReview(ctx context.Context, professionalId int64, page int, pageSize int) (Reviews, error)
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
	//
	// 23005 means unique_violation, defined in pq.Error
	//
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return utils.ErrDuplicate
	}
	return err
}

func (s *serviceImpl) Update(ctx context.Context, id string, p UpdateRequest) error {
	return s.repository.Update(ctx, id, p)
}

func (s *serviceImpl) FindAllReview(ctx context.Context, professionalId int64, page int, pageSize int) (Reviews, error) {
	return s.repository.FindAllReview(ctx, professionalId, page, pageSize)
}
