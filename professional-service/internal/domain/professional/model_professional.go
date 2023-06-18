package professional

import "time"

type Professional struct {
	ID              int        `json:"id,omitempty"`
	Email           string     `json:"email,omitempty"`
	Password        string     `json:"password,omitempty"`
	FirstName       *string     `json:"first_name,omitempty"`
	LastName        *string     `json:"last_name,omitempty"`
	CoachType       *string     `json:"coach_type,omitempty"`
	PriceNumber     *int        `json:"price_number,omitempty"`
	PriceCurrency   *string     `json:"price_currency,omitempty"`
	ProfileImageUrl *string     `json:"profile_image_url,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}
