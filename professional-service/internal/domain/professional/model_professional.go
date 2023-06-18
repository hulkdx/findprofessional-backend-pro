package professional

import "time"

type Professional struct {
	ID              int        `json:"id,omitempty"`
	Email           string     `json:"email,omitempty"`
	Password        string     `json:"password,omitempty"`
	FirstName       *string    `json:"firstName,omitempty"`
	LastName        *string    `json:"lastName,omitempty"`
	CoachType       *string    `json:"coachType,omitempty"`
	PriceNumber     *int       `json:"priceNumber,omitempty"`
	PriceCurrency   *string    `json:"priceCurrency,omitempty"`
	ProfileImageUrl *string    `json:"profileImageUrl,omitempty"`
	Rating          *string    `json:"rating,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
}
