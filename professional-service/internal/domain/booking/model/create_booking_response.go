package bookingmodel

type PaymentResponse struct {
	PaymentIntent               string `json:"payment_intent"`
	Customer                    string `json:"customer"`
	PublishableKey              string `json:"publishable_key"`
	CustomerSessionClientSecret string `json:"customer_session_client_secret"`
}
