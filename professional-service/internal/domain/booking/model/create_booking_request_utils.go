package booking_model

import (
	"net/http"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

func ParseCreateRequest(r *http.Request) (*CreateBookingRequest, error) {
	request := CreateBookingRequest{}
	if err := utils.ReadJSON(r, &request); err != nil {
		return nil, err
	}
	if err := utils.Validate(&request); err != nil {
		return nil, err
	}
	return &request, nil
}
