package professional

import (
	"database/sql"
)

type Repository interface {
	FindAll() ([]Professional, error)
}

func NewRepository(db *sql.DB) Repository {
	return &impl{db}
}

type impl struct {
	db *sql.DB
}

func (p *impl) FindAll() ([]Professional, error) {
	panic("TODO")
}
