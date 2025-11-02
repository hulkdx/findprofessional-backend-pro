package professional

import (
	"context"
	"database/sql"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindAll(ctx context.Context) ([]model_professional.Professional, error)
	FindById(ctx context.Context, id string) (model_professional.Professional, error)
	Create(ctx context.Context, request CreateRequest, pending bool) error
	Update(ctx context.Context, id string, p UpdateRequest) error
	FindAllReview(ctx context.Context, professionalId int64, page int, pageSize int) (model_professional.Reviews, error)
	GetAvailability(ctx context.Context, professionalId int64) (model_professional.Availabilities, error)
	UpdateAvailability(ctx context.Context, professionalId int64, availability UpdateAvailabilityRequest) error
}

type repositoryImpl struct {
	db           *pgxpool.Pool
	timeProvider utils.TimeProvider
}

func NewRepository(db *pgxpool.Pool, timeProvider utils.TimeProvider) Repository {
	return &repositoryImpl{
		db,
		timeProvider,
	}
}

func (r *repositoryImpl) Update(ctx context.Context, id string, p UpdateRequest) error {
	query := "UPDATE professionals SET updated_at = $2, first_name = $3, last_name = $4, coach_type = $5"
	args := []any{
		id,
		time.Now(),
		p.FirstName,
		p.LastName,
		p.CoachType,
	}
	if p.Email != nil {
		query += ", email = $6"
		args = append(args, *p.Email)
	}
	if p.Price != nil && p.PriceCurrency != nil {
		query += ", price_number = $7, price_currency = $8"
		args = append(args, p.Price, p.PriceCurrency)
	}
	if p.ProfileImageUrl != nil {
		query += ", profile_image_url = $9"
		args = append(args, *p.ProfileImageUrl)
	}
	if p.Description != nil {
		query += ", description = $10"
		args = append(args, *p.Description)
	}
	if p.SkypeId != nil {
		query += ", skype_id = $11"
		args = append(args, *p.SkypeId)
	}

	query += " WHERE id = $1"
	return performUpdate(r.db, ctx, query, args...)
}

func (r *repositoryImpl) FindAllReview(ctx context.Context, professionalID int64, page int, pageSize int) (model_professional.Reviews, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT 
			pr.id,
			pr.rate,
			pr.content_text,
			pr.created_at,
			pr.updated_at,
			u.id,
			u.email,
			u.first_name,
			u.last_name,
			u.profile_image
		FROM professional_review pr
		LEFT JOIN users u ON pr.user_id = u.id
		WHERE pr.professional_id = $1
		ORDER BY pr.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, professionalID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make(model_professional.Reviews, 0)

	for rows.Next() {
		var review model_professional.Review

		err := rows.Scan(
			&review.ID,
			&review.Rate,
			&review.ContentText,
			&review.CreatedAt,
			&review.UpdatedAt,
			&review.User.ID,
			&review.User.Email,
			&review.User.FirstName,
			&review.User.LastName,
			&review.User.ProfileImage,
		)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (r *repositoryImpl) Create(ctx context.Context, request CreateRequest, pending bool) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
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
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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

	return tx.Commit(ctx)
}

func (r *repositoryImpl) GetAvailability(ctx context.Context, professionalId int64) (model_professional.Availabilities, error) {
	query := `
		SELECT
			id,
			LOWER(availability) AT TIME ZONE 'UTC',
			UPPER(availability) AT TIME ZONE 'UTC',
			created_at,
			updated_at
	FROM professional_availability
	WHERE
		professional_id = $1 AND
		LOWER(availability) > $2`
	rows, err := r.db.Query(ctx, query,
		professionalId,
		r.timeProvider.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	availabilities := make(model_professional.Availabilities, 0)
	for rows.Next() {
		var availability model_professional.Availability
		var from time.Time
		var to time.Time

		err := rows.Scan(
			&availability.ID,
			&from,
			&to,
			&availability.CreatedAt,
			&availability.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		availability.From = from
		availability.To = to

		availabilities = append(availabilities, availability)
	}

	return availabilities, nil
}

func (r *repositoryImpl) UpdateAvailability(ctx context.Context, professionalId int64, availability UpdateAvailabilityRequest) error {
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
