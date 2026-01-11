package professionalrepo

import (
	"context"
	"fmt"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
)

func (r *RepositoryImpl) GetBookings(ctx context.Context, id int64, idType professional.UserType) (*professional.Bookings, error) {
	const query = `
	SELECT
		b.id,
		b.status,
		b.scheduled_start_at,
		b.scheduled_end_at,
		b.total_amount_cents,
		b.currency,
		b.created_at,
		b.updated_at,

		p.id,
		p.first_name,
		p.last_name,
		COALESCE(NULLIF(b.session_platform, ''), NULLIF(p.session_platform, '')),
		COALESCE(NULLIF(b.session_link, ''),     NULLIF(p.session_link, '')),

		u.id,
		u.first_name,
		u.last_name,
		u.email
	FROM
		bookings b
	JOIN professionals p ON p.id = b.professional_id
	JOIN users u ON u.id = b.user_id

	WHERE %s = $1
	ORDER BY b.scheduled_start_at DESC
	`
	var finalQuery string
	switch idType {
	case professional.UserTypePro:
		finalQuery = fmt.Sprintf(query, "b.professional_id")
	case professional.UserTypeNormal:
		finalQuery = fmt.Sprintf(query, "b.user_id")
	}

	rows, err := r.db.Query(ctx, finalQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make(professional.Bookings, 0)
	for rows.Next() {
		var item professional.Booking
		err := rows.Scan(
			&item.ID,
			&item.Status,
			&item.ScheduledStartAt,
			&item.ScheduledEndAt,
			&item.TotalAmountCents,
			&item.Currency,
			&item.CreatedAt,
			&item.UpdatedAt,

			&item.Professional.ID,
			&item.Professional.FirstName,
			&item.Professional.LastName,
			&item.Session.Platform,
			&item.Session.Link,

			&item.User.ID,
			&item.User.FirstName,
			&item.User.LastName,
			&item.User.Email,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return &items, nil
}
