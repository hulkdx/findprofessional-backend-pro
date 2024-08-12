package professional

import (
	"context"
	"database/sql"
)

func performUpdateTx(tx *sql.Tx, ctx context.Context, query string, args ...any) error {
	return performUpdateCheckErrors(tx.ExecContext(ctx, query, args))
}

func performUpdate(db *sql.DB, ctx context.Context, query string, args ...any) error {
	return performUpdateCheckErrors(db.ExecContext(ctx, query, args))
}

func performUpdateCheckErrors(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
