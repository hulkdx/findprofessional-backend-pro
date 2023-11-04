package professional

import (
	"encoding/json"
	"time"

	"cloud.google.com/go/civil"
)

// TODO: convert id to int64
type Professional struct {
	ID              int            `json:"id,omitempty"`
	Email           string         `json:"email,omitempty"`
	Password        string         `json:"password,omitempty"`
	FirstName       string         `json:"firstName,omitempty"`
	LastName        string         `json:"lastName,omitempty"`
	CoachType       *string        `json:"coachType,omitempty"`
	PriceNumber     *int           `json:"priceNumber,omitempty"`
	PriceCurrency   *string        `json:"priceCurrency,omitempty"`
	ProfileImageUrl *string        `json:"profileImageUrl,omitempty"`
	Description     *string        `json:"description,omitempty"`
	Rating          *string        `json:"rating,omitempty"`
	Availability    Availabilities `json:"availabilities,omitempty"`
	CreatedAt       *time.Time     `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time     `json:"updatedAt,omitempty"`
}

// TODO: remove gorm
type Availability struct {
	ID             int
	ProfessionalID int
	Date           civil.Date `json:"date,omitempty" gorm:"type:date;serializer:json"`
	From           civil.Time `json:"from,omitempty" gorm:"serializer:custom"`
	To             civil.Time `json:"to,omitempty" gorm:"serializer:custom"`
}

type Availabilities []Availability

func (ls *Availabilities) Scan(src any) error {
	if src == nil {
		return nil
	}
	return json.Unmarshal(src.([]byte), ls)
}
