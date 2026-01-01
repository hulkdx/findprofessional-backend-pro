package professionalrepo

import (
	"context"
	"database/sql"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
)

func (r *RepositoryImpl) Create(ctx context.Context, request professional.CreateRequest, pending bool) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	txDone := false
	defer func() {
		if !txDone {
			tx.Rollback(ctx)
		}
	}()

	query := `
		INSERT INTO professionals (
			email,
			password,
			first_name,
			last_name,
			coach_type,
			description,
			price_number,
			price_currency,
			pending,
			session_link,
			session_platform,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id;
	`

	var professionalId int64
	row := tx.QueryRow(ctx, query,
		request.Email,
		request.Password,
		request.FirstName,
		request.LastName,
		request.CoachType,
		request.AboutMe,
		request.Price,
		request.PriceCurrency,
		pending,
		request.SessionLink,
		request.SessionPlatform,
		r.timeProvider.Now(),
		r.timeProvider.Now(),
	)
	err = row.Scan(&professionalId)
	if err != nil {
		return err
	}

	query = "UPDATE users SET professional_id = $1, updated_at = $2 WHERE email = $3"
	result, err := tx.Exec(ctx, query,
		professionalId,
		r.timeProvider.Now(),
		request.Email,
	)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	txDone = true
	return nil
}
