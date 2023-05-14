package db

import (
	"database/sql"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/config"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/logger"

	_ "github.com/lib/pq"
)

func Connect(cfg config.DatabaseConfig) *sql.DB {
	db, err := sql.Open("postgres", cfg.Dsn)
	if err != nil {
		defer Close(db)
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		defer Close(db)
		panic(err)
	}
	return db
}

func Close(db *sql.DB) {
	err := db.Close()
	if err != nil {
		logger.Error("Database close error:", err)
	}
}
