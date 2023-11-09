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
	Review          Reviews        `json:"reviews,omitempty"`
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

type Reviews []Review

func (ls *Reviews) Scan(src any) error {
	if src == nil {
		return nil
	}
	return json.Unmarshal(src.([]byte), ls)
}

type Review struct {
	ID             uint      `json:"-"`
	UserID         int64     `json:"-"`
	ProfessionalID int64     `json:"-"`
	Rate           int       `json:"rate,omitempty"`
	ContentText    string    `json:"content_text,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
}
