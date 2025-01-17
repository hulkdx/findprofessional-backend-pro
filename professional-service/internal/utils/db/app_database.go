package db

import (
	"context"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/config"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"
)

func Connect(ctx context.Context, cfg config.DatabaseConfig) *pgxpool.Pool {
	db, err := pgxpool.New(ctx, cfg.Dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping(ctx)
	if err != nil {
		defer Close(db)
		panic(err)
	}

	return db
}

func Close(db *pgxpool.Pool) {
	db.Close()
}
