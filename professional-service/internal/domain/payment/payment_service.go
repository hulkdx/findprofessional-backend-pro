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
	CreatePaymentIntent(
		ctx context.Context,
		userId int64,
		req *booking_model.CreateBookingRequest,
		auth string,
	) (*booking_model.PaymentIntentResponse, error)
}

type paymentServiceImpl struct {
	httpClient *http.Client
}

func NewService() PaymentService {
	return &paymentServiceImpl{
		httpClient: utils.CreateDefaultAppHttpClient(),
	}
}

func (s *paymentServiceImpl) CreatePaymentIntent(
	ctx context.Context,
	userId int64,
	req *booking_model.CreateBookingRequest,
	auth string,
) (*booking_model.PaymentIntentResponse, error) {
	url := fmt.Sprintf("%s/payments/create-intent", baseUrl)
	request := &PaymentRequest{
		AmountsInCents: req.AmountInCents,
		Currency:       req.Currency,
	}
	var response booking_model.PaymentIntentResponse
	requestHeader := &http.Header{}
	requestHeader.Set("Authorization", auth)
	requestHeader.Set("Idempotency-Key", req.IdempotencyKey)
	err := utils.DoHttpRequestAsStruct(ctx, s.httpClient, http.MethodPost, url, &request, &response, requestHeader)
	return &response, err
}
