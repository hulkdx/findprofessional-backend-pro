package booking_model

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
