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
	query := "SELECT id, email, password, created_at, updated_at FROM professionals"
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	professionals := []Professional{}
	for rows.Next() {
		var pro Professional
		err := rows.Scan(&pro.ID, &pro.Email, &pro.Password, &pro.CreatedAt, &pro.UpdatedAt)
		if err != nil {
			return nil, err
		}
		professionals = append(professionals, pro)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return professionals, nil
}
