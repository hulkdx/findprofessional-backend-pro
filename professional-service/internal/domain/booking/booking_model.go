package booking

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
	BookingID     string `json:"booking_id"`
	AmountInCents int64  `json:"amount_in_cents"`
	Currency      string `json:"currency"`
	ClientSecret  string `json:"client_secret"`
	CustomerID    string `json:"customer_id"`
	EphemeralKey  string `json:"ephemeral_key"`
	HoldExpiresAt string `json:"hold_expires_at"`
}

// TODO: create a enum class for booking status
type BookingStatus string

const (
	BookingStatusHold      BookingStatus = "hold"
	BookingStatusCompleted               = "completed"
	BookingStatusCanceled                = "canceled"
)
