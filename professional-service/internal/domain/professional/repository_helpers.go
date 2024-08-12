package professional

import (
	"context"
	"database/sql"
)

func (r *repositoryImpl) performUpdate(ctx context.Context, query string, args ...any) error {
	result, err := r.db.ExecContext(ctx, query, args)
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
