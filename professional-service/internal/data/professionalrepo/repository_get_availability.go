package professionalrepo

import (
	"context"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
)

func (r *RepositoryImpl) GetAvailability(ctx context.Context, professionalId int64) (professional.Availabilities, error) {
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
	availabilities := make(professional.Availabilities, 0)
	for rows.Next() {
		var availability professional.Availability
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
