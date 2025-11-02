package professionalrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/jackc/pgx/v5"
)

func (r *RepositoryImpl) UpdateAvailability(ctx context.Context, professionalId int64, availability professional.UpdateAvailabilityRequest) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	txDone := false
	defer func() {
		if txDone {
			tx.Commit(ctx)
		} else {
			tx.Rollback(ctx)
		}
	}()

	now := r.timeProvider.Now()
	rows := make([][]any, len(availability.Items))

	var min *time.Time
	var max *time.Time
	for i, e := range availability.Items {
		if i == 0 {
			min = &e.From
			max = &e.To
		} else {
			min = utils.MinTime(min, &e.From)
			max = utils.MaxTime(max, &e.To)
		}

		rows[i] = []any{
			professionalId,
			utils.ConvertToTstzrange(e.From, e.To),
			now,
			now,
		}
	}

	if len(availability.Items) > 0 {
		query := `
			DELETE FROM professional_availability
			WHERE
				professional_id = $1 AND
				availability && tstzrange($2, $3);
		`
		_, err = tx.Exec(ctx, query, professionalId, min, max)
		if err != nil {
			return err
		}
	}

	columns := []string{
		"professional_id",
		"availability",
		"created_at",
		"updated_at",
	}
	count, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"professional_availability"},
		columns,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}
	if count != int64(len(rows)) {
		return sql.ErrNoRows
	}

	txDone = true
	return nil
}
