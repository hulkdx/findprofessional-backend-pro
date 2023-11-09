package professional

import (
	"encoding/json"
	"time"

	"cloud.google.com/go/civil"
)

type Professional struct {
	ID              int64          `json:"id,omitempty"`
	Email           string         `json:"email,omitempty"`
	Password        string         `json:"password,omitempty"`
	FirstName       string         `json:"firstName,omitempty"`
	LastName        string         `json:"lastName,omitempty"`
	CoachType       string         `json:"coachType,omitempty"`
	PriceNumber     *int           `json:"priceNumber,omitempty"`
	PriceCurrency   *string        `json:"priceCurrency,omitempty"`
	ProfileImageUrl *string        `json:"profileImageUrl,omitempty"`
	Description     *string        `json:"description,omitempty"`
	Rating          *string        `json:"rating,omitempty"`
	Availability    Availabilities `json:"availability,omitempty"`
	Review          Review         `json:"review"`
	CreatedAt       time.Time      `json:"createdAt,omitempty"`
	UpdatedAt       time.Time      `json:"updatedAt,omitempty"`
}

type Availability struct {
	ID             int64      `json:"-"`
	ProfessionalID int        `json:"-"`
	Date           civil.Date `json:"date,omitempty"`
	From           civil.Time `json:"from,omitempty"`
	To             civil.Time `json:"to,omitempty"`
	CreatedAt      time.Time  `json:"createdAt,omitempty"`
	UpdatedAt      time.Time  `json:"updatedAt,omitempty"`
}

type Availabilities []Availability

func (ls *Availabilities) Scan(src any) error {
	if src == nil {
		return nil
	}
	return json.Unmarshal(src.([]byte), ls)
}

type Review struct {
	Total   int                  `json:"total"`
	Content []ProfessionalRating `json:"content,omitempty"`
}

type ProfessionalRating struct {
	ID             uint
	UserID         int64
	ProfessionalID int64
	Rate           int
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
}
