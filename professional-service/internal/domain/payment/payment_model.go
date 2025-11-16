package payment

import bookingmodel "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"

type PaymentRequest struct {
	AmountsInCents int64                       `json:"amounts_in_cents"`
	Currency       string                      `json:"currency"`
	HoldId         int64                       `json:"hold_id"`
	ProfessionalId int64                       `json:"professional_id"`
	Availabilities []bookingmodel.Availability `json:"availabilities" validate:"required"`
}

type PaymentResponse struct {
}
