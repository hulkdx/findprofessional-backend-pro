package bookingmodel

type CreateBookingRequest struct {
	Availabilities   []Availability `json:"availabilities" validate:"required"`
	IdempotencyKey   string         `json:"idempotency_key"  validate:"required,max=36,min=32"`
	AmountInCents    int64          `json:"amount_in_cents"  validate:"required"`
	Currency         string         `json:"currency"  validate:"required,len=3"`
	StripeApiVersion string         `json:"stripe_api_version" validate:"required,max=50"`
}

type Availability struct {
	Id int64 `json:"id" validate:"required"`
}
