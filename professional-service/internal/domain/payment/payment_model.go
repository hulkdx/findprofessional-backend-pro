package payment

type PaymentRequest struct {
	AmountsInCents int64  `json:"amounts_in_cents"`
	Currency       string `json:"currency"`
	HoldId         int64  `json:"hold_id"`
}

type PaymentResponse struct {
}
