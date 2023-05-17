package professional

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Repository interface {
	FindAll(fields ...string) ([]Professional, error)
}

type repositoryImpl struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repositoryImpl{db}
}

func (p *repositoryImpl) FindAll(fields ...string) ([]Professional, error) {
	selectedFields := strings.Join(fields, ", ")
	query := fmt.Sprintf("SELECT %s FROM professionals", selectedFields)
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	professionals := []Professional{}
	for rows.Next() {
		pro := Professional{}
		elem := reflect.ValueOf(&pro).Elem()
		scanArgs := make([]interface{}, len(fields))
		for i := range fields {
			field := elem.FieldByName(fields[i])
			scanArgs[i] = field.Addr().Interface()
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}
		professionals = append(professionals, pro)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return professionals, nil
}
