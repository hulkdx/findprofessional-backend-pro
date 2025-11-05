package bookingmodel

import "time"

type BookingHold struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	IdempotencyKey string    `json:"idempotency_key"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}
