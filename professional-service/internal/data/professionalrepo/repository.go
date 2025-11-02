package professionalrepo

import (
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryImpl struct {
	db           *pgxpool.Pool
	timeProvider utils.TimeProvider
}

func NewRepository(db *pgxpool.Pool, timeProvider utils.TimeProvider) professional.Repository {
	return &RepositoryImpl{
		db,
		timeProvider,
	}
}
