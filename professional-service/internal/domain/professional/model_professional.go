package professional

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/user"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
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
	ReviewSize      int64          `json:"reviewSize"`
	Rating          *string        `json:"rating,omitempty"`
	Availability    Availabilities `json:"availability,omitempty"`
	Review          Reviews        `json:"reviews,omitempty"`
	Pending         bool           `json:"-"`
	CreatedAt       time.Time      `json:"createdAt,omitempty"`
	UpdatedAt       time.Time      `json:"updatedAt,omitempty"`
}

type Availability struct {
	ID             int64      `json:"-"`
	ProfessionalID int64      `json:"-"`
	Date           civil.Date `json:"date,omitempty"`
	From           civil.Time `json:"from,omitempty"`
	To             civil.Time `json:"to,omitempty"`
	CreatedAt      time.Time  `json:"createdAt,omitempty"`
	UpdatedAt      time.Time  `json:"updatedAt,omitempty"`
}

type Availabilities []Availability

func (ls *Availabilities) Scan(src any) error {
	return utils.Unmarshal(src, ls)
}

type Reviews []Review

func (ls *Reviews) Scan(src any) error {
	return utils.Unmarshal(src, ls)
}

type Review struct {
	ID             int64     `json:"id,omitempty"`
	UserID         int64     `json:"-"`
	User           user.User `json:"user,omitempty"`
	ProfessionalID int64     `json:"-"`
	Rate           int       `json:"rate,omitempty"`
	ContentText    *string   `json:"contentText,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
}
