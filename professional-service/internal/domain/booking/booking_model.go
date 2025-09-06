package booking

import "time"

type Slot struct {
	Start time.Time `json:"start"  validate:"required"`
	End   time.Time `json:"end"  validate:"required"`
}

type CreateBookingRequest struct {
	Slots          []Slot `json:"slots" validate:"required,max=50"`
	IdempotencyKey string `json:"idempotency_key"  validate:"required,max=50"`
	AmountInCents  int64  `json:"amount_in_cents"  validate:"required"`
	Currency       string `json:"currency"  validate:"required,len=3"`
}

type CreateBookingResponse struct {
	BookingID     string `json:"booking_id"`
	AmountInCents int64  `json:"amount_in_cents"`
	Currency      string `json:"currency"`
	ClientSecret  string `json:"client_secret"`
	CustomerID    string `json:"customer_id"`
	EphemeralKey  string `json:"ephemeral_key"`
	HoldExpiresAt string `json:"hold_expires_at"`
}
