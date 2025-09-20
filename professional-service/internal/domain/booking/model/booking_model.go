package booking_model

import "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/payment"

type Slot struct {
	Date string `json:"date" validate:"required,max=50"`
	From string `json:"from" validate:"required,max=50"`
	To   string `json:"to" validate:"required,max=50"`
}

type CreateBookingRequest struct {
	Slots          []Slot `json:"slots" validate:"required,max=50"`
	IdempotencyKey string `json:"idempotency_key"  validate:"required,max=50"`
	AmountInCents  int64  `json:"amount_in_cents"  validate:"required"`
	Currency       string `json:"currency"  validate:"required,len=3"`
}

type CreateBookingResponse struct {
	BookingID             int64                         `json:"booking_id"`
	PaymentIntentResponse payment.PaymentIntentResponse `json:"payment_intent_response"`
}

type BookingStatus string

const (
	BookingStatusHold      BookingStatus = "hold"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusCanceled  BookingStatus = "canceled"
)