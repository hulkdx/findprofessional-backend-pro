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
		idempotencyKey string,
		auth string,
		professionalId int64,
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
	idempotencyKey string,
	auth string,
	professionalId int64,
) (*bookingmodel.PaymentIntentResponse, error) {
	url := fmt.Sprintf("%s/payments/create-intent", baseUrl)
	request := &PaymentRequest{
		AmountsInCents: AmountInCents,
		Currency:       Currency,
		HoldId:         holdId,
		ProfessionalId: professionalId,
	}
	var response bookingmodel.PaymentIntentResponse
	requestHeader := &http.Header{}
	requestHeader.Set("Authorization", auth)
	requestHeader.Set("Idempotency-Key", idempotencyKey)
	err := utils.DoHttpRequestAsStruct(ctx, s.httpClient, http.MethodPost, url, &request, &response, requestHeader)
	return &response, err
}
