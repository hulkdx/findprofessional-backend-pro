package professional

import "time"

type UserType string

const (
	UserTypeNormal UserType = "user"
	UserTypePro    UserType = "professional"
)

type Party struct {
	ID        int64   `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     *string `json:"email,omitempty"`
}

type SessionInfo struct {
	Platform *string `json:"sessionPlatform,omitempty"`
	Link     *string `json:"sessionLink,omitempty"`
}

type Booking struct {
	ID               int64     `json:"id"`
	Status           string    `json:"status"`
	ScheduledStartAt time.Time `json:"scheduledStartAt"`
	ScheduledEndAt   time.Time `json:"scheduledEndAt"`
	TotalAmountCents int64     `json:"totalAmountCents"`
	Currency         string    `json:"currency"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`

	Professional Party       `json:"professional"`
	User         Party       `json:"user"`
	Session      SessionInfo `json:"session"`
}

type Bookings []Booking
