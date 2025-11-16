package payment

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

const baseUrl = "http://payment-service:8082"

type Service interface {
	CreatePaymentIntent(
		ctx context.Context,
		request PaymentRequest,
		idempotencyKey string,
		auth string,
	) (*bookingmodel.PaymentResponse, error)
}

type paymentServiceImpl struct {
	httpClient *http.Client
}

func NewService() Service {
	return &paymentServiceImpl{
		httpClient: utils.CreateDefaultAppHttpClient(),
	}
}

func (s *paymentServiceImpl) CreatePaymentIntent(
	ctx context.Context,
	request PaymentRequest,
	idempotencyKey string,
	auth string,
) (*bookingmodel.PaymentResponse, error) {
	url := fmt.Sprintf("%s/payments/create-intent", baseUrl)
	var response bookingmodel.PaymentResponse
	requestHeader := &http.Header{}
	requestHeader.Set("Authorization", auth)
	requestHeader.Set("Idempotency-Key", idempotencyKey)
	err := utils.DoHttpRequestAsStruct(ctx, s.httpClient, http.MethodPost, url, request, &response, requestHeader)
	return &response, err
}
