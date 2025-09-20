package booking_model

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
	BookingID             int64                 `json:"booking_id"`
	PaymentIntentResponse PaymentIntentResponse `json:"payment_intent_response"`
}

type PaymentIntentResponse struct {
	PaymentIntent  string `json:"payment_intent"`
	EphemeralKey   string `json:"ephemeral_key"`
	Customer       string `json:"customer"`
	PublishableKey string `json:"publishable_key"`
}

type BookingStatus string

const (
	BookingStatusHold      BookingStatus = "hold"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusCanceled  BookingStatus = "canceled"
)
