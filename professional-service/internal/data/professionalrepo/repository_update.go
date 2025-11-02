package professionalrepo

import (
	"context"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/data/utils/sqlutils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
)

func (r *RepositoryImpl) Update(ctx context.Context, id string, p professional.UpdateRequest) error {
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
	return sqlutils.PerformUpdate(r.db, ctx, query, args...)
}
