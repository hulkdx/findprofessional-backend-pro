package db

import (
	"database/sql"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/config"
)

func Connect(cfg config.DatabaseConfig) *sql.DB {
	db, err := sql.Open("postgres", cfg.Dsn())
	if err != nil {
		db.Close()
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		panic(err)
	}
	return db
}
