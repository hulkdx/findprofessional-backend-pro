package payment

type PaymentRequest struct {
	AmountsInCents int64  `json:"amounts_in_cents" validate:"required"`
	Currency       string `json:"currency" validate:"required,max=50"`
}
