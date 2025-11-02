package sqlutils

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PerformUpdateTx(tx pgx.Tx, ctx context.Context, query string, args ...any) error {
	return performUpdateCheckErrors(tx.Exec(ctx, query, args...))
}

func PerformUpdate(db *pgxpool.Pool, ctx context.Context, query string, args ...any) error {
	return performUpdateCheckErrors(db.Exec(ctx, query, args...))
}

func performUpdateCheckErrors(result pgconn.CommandTag, err error) error {
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
