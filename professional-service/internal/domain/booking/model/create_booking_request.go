package booking_model

type CreateBookingRequest struct {
	Slots          []Slot `json:"slots" validate:"required,max=50"`
	IdempotencyKey string `json:"idempotency_key"  validate:"required,max=50"`
	AmountInCents  int64  `json:"amount_in_cents"  validate:"required"`
	Currency       string `json:"currency"  validate:"required,len=3"`
}

type Slot struct {
	Date string `json:"date" validate:"required,max=50"`
	From string `json:"from" validate:"required,max=50"`
	To   string `json:"to" validate:"required,max=50"`
}
