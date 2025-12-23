package professionalrepo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/data/utils/sqlutils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
)

func (r *RepositoryImpl) Update(ctx context.Context, id string, p professional.UpdateRequest) error {
	if (p.Price == nil) != (p.PriceCurrency == nil) {
		return fmt.Errorf("price and price_currency must be provided together")
	}

	args := []any{}
	set := []string{}

	args = append(args, id)

	add := func(col string, val any) {
		set = append(set, fmt.Sprintf("%s = $%d", col, len(args)+1))
		args = append(args, val)
	}

	add("updated_at", time.Now().UTC())
	add("first_name", p.FirstName)
	add("last_name", p.LastName)
	add("coach_type", p.CoachType)

	if p.Email != nil {
		add("email", *p.Email)
	}
	if p.Price != nil && p.PriceCurrency != nil {
		add("price_number", *p.Price)
		add("price_currency", *p.PriceCurrency)
	}
	if p.ProfileImageUrl != nil {
		add("profile_image_url", *p.ProfileImageUrl)
	}
	if p.Description != nil {
		add("description", *p.Description)
	}
	if p.SkypeId != nil {
		add("skype_id", *p.SkypeId)
	}

	query := fmt.Sprintf("UPDATE professionals SET %s WHERE id = $1", strings.Join(set, ", "))
	return sqlutils.PerformUpdate(r.db, ctx, query, args...)
}
