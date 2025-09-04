package booking

import "time"

type Slot struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type CreateBookingRequest struct {
	Slots          []Slot `json:"slots"`
	IdempotencyKey string `json:"idempotency_key"`
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
