package professional

import (
	"encoding/json"
	"time"
)

type Professional struct {
	ID              int            `json:"id,omitempty"`
	Email           string         `json:"email,omitempty"`
	Password        string         `json:"password,omitempty"`
	FirstName       string         `json:"firstName,omitempty"`
	LastName        string         `json:"lastName,omitempty"`
	CoachType       string         `json:"coachType,omitempty"`
	PriceNumber     int            `json:"priceNumber,omitempty"`
	PriceCurrency   string         `json:"priceCurrency,omitempty"`
	ProfileImageUrl *string        `json:"profileImageUrl,omitempty"`
	Description     *string        `json:"description,omitempty"`
	Rating          *string        `json:"rating,omitempty"`
	Availability    Availabilities `json:"availabilities,omitempty"`
	CreatedAt       *time.Time     `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time     `json:"updatedAt,omitempty"`
}

type Availability struct {
	ID             int
	ProfessionalID int
	Date           time.Time `json:"date"`
	From           time.Time `json:"from"`
	To             time.Time `json:"to"`
}

type Availabilities []Availability

func (ls *Availabilities) Scan(src any) error {
	return json.Unmarshal(src.([]byte), ls)
}
