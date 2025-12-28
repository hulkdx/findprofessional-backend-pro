package professional

import (
	"net/http"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

type CreateRequest struct {
	Email         string `json:"email" validate:"email,required,max=50"`
	Password      string `json:"password" validate:"max=50"`
	FirstName     string `json:"firstName" validate:"max=50"`
	LastName      string `json:"lastName" validate:"max=50"`
	AboutMe       string `json:"aboutMe" validate:"max=500"`
	CoachType     string `json:"coachType" validate:"max=50"`
	Price         int64  `json:"price"`
	PriceCurrency string `json:"priceCurrency" validate:"max=50"`
}

func parseCreateRequest(r *http.Request) (CreateRequest, error) {
	request := CreateRequest{}
	if err := utils.ReadJSON(r, &request); err != nil {
		return CreateRequest{}, err
	}
	if err := utils.Validate(&request); err != nil {
		return CreateRequest{}, err
	}
	return request, nil
}
