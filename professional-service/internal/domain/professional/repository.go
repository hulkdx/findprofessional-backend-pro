package professional

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Repository interface {
	FindAll(fields ...string) ([]Professional, error)
}

func NewRepository(db *sql.DB) Repository {
	return &impl{db}
}

type impl struct {
	db *sql.DB
}

func (p *impl) FindAll(fields ...string) ([]Professional, error) {
	selectedFields := strings.Join(fields, ", ")
	query := fmt.Sprintf("SELECT %s FROM professionals", selectedFields)
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	professionals := []Professional{}
	for rows.Next() {
		scanArgs := make([]any, len(fields))
		for i := range fields {
			scanArgs[i] = new(any)
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		pro := Professional{}
		for i, field := range fields {
			item := scanArgs[i].(*interface{})
			switch field {
			case "id":
				pro.ID = int((*item).(int64))
			case "email":
				pro.Email = (*item).(string)
			case "password":
				pro.Password = (*item).(string)
			case "created_at":
				time := (*item).(time.Time)
				pro.CreatedAt = &time
			case "updated_at":
				time := (*item).(time.Time)
				pro.UpdatedAt = &time
			}
		}
		professionals = append(professionals, pro)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return professionals, nil
}
