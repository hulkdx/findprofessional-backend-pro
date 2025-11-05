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
		holdId int64,
		AmountInCents int64,
		Currency string,
		IdempotencyKey string,
		auth string,
	) (*bookingmodel.PaymentIntentResponse, error)
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
	holdId int64,
	AmountInCents int64,
	Currency string,
	IdempotencyKey string,
	auth string,
) (*bookingmodel.PaymentIntentResponse, error) {
	url := fmt.Sprintf("%s/payments/create-intent", baseUrl)
	request := &PaymentRequest{
		AmountsInCents: AmountInCents,
		Currency:       Currency,
		HoldId:         holdId,
	}
	var response bookingmodel.PaymentIntentResponse
	requestHeader := &http.Header{}
	requestHeader.Set("Authorization", auth)
	requestHeader.Set("Idempotency-Key", IdempotencyKey)
	err := utils.DoHttpRequestAsStruct(ctx, s.httpClient, http.MethodPost, url, &request, &response, requestHeader)
	return &response, err
}
