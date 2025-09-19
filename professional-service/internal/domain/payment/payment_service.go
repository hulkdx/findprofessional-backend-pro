package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

const baseUrl = "http://payment-service:8082"

type PaymentService interface {
	CreatePaymentIntent(ctx context.Context, userId int64, amountInCents int64, currency string) (*PaymentIntentResponse, error)
}

type paymentServiceImpl struct {
	httpClient *http.Client
}

func NewService() PaymentService {
	return &paymentServiceImpl{
		httpClient: utils.CreateDefaultAppHttpClient(),
	}
}

type PaymentIntentResponse struct {
	PaymentIntent  string `json:"payment_intent"`
	EphemeralKey   string `json:"ephemeral_key"`
	Customer       string `json:"customer"`
	PublishableKey string `json:"publishable_key"`
}

func (s *paymentServiceImpl) CreatePaymentIntent(ctx context.Context, userId int64, amountInCents int64, currency string) (*PaymentIntentResponse, error) {
	url := fmt.Sprintf("%s/payments/create-intent", baseUrl)
	request := fmt.Sprintf(`{"amountInCents": %d, "currency": "%s"}`, amountInCents, currency)
	bodyResponse, err := utils.DoHttpRequestAsReader(ctx, s.httpClient, http.MethodPost, url, request)
	if err != nil {
		return nil, err
	}
	defer bodyResponse.Close()

	var paymentIntentResponse PaymentIntentResponse
	if err := json.NewDecoder(bodyResponse).Decode(&paymentIntentResponse); err != nil {
		return nil, err
	}
	return &paymentIntentResponse, nil
}
