package booking

import "github.com/jackc/pgx/v5/pgxpool"

type Repository interface {
}

type repositoryImpl struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repositoryImpl{
		db: db,
	}
}
