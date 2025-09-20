package payment

import (
	"context"
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
	request := &PaymentRequest{
		AmountsInCents: amountInCents,
		Currency:       currency,
	}
	var response booking_model.PaymentIntentResponse
	err := utils.DoHttpRequestAsStruct(ctx, s.httpClient, http.MethodPost, url, &request, &response)
	return &response, err
}
