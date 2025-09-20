package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

const baseUrl = "http://payment-service:8082"

type PaymentService interface {
	CreatePaymentIntent(ctx context.Context, userId int64, amountInCents int64, currency string) (*booking_model.PaymentIntentResponse, error)
}

type paymentServiceImpl struct {
	httpClient *http.Client
}

func NewService() PaymentService {
	return &paymentServiceImpl{
		httpClient: utils.CreateDefaultAppHttpClient(),
	}
}

func (s *paymentServiceImpl) CreatePaymentIntent(ctx context.Context, userId int64, amountInCents int64, currency string) (*booking_model.PaymentIntentResponse, error) {
	url := fmt.Sprintf("%s/payments/create-intent", baseUrl)
	request := fmt.Sprintf(`{"amountInCents": %d, "currency": "%s"}`, amountInCents, currency)
	bodyResponse, err := utils.DoHttpRequestAsReader(ctx, s.httpClient, http.MethodPost, url, request)
	if err != nil {
		return nil, err
	}
	defer bodyResponse.Close()

	var paymentIntentResponse booking_model.PaymentIntentResponse
	if err := json.NewDecoder(bodyResponse).Decode(&paymentIntentResponse); err != nil {
		return nil, err
	}
	return &paymentIntentResponse, nil
}
