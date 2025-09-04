package booking

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

func (c *BookingController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth := r.Header.Get("Authorization")
	if !c.userService.IsAuthenticated(ctx, auth) {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	proID := chi.URLParam(r, "id")
	createBookingRequest, err := parseCreateRequest(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	booking, err := c.service.Create(ctx, proID, createBookingRequest)
	if err != nil {
		utils.WriteGeneralError(w, utils.ErrUnknown)
		return
	}
	fmt.Printf("Booking created: %+v\n", booking)
}

func parseCreateRequest(r *http.Request) (CreateBookingRequest, error) {
	request := CreateBookingRequest{}
	if err := utils.ReadJSON(r, &request); err != nil {
		return CreateBookingRequest{}, err
	}
	if err := utils.Validate(request); err != nil {
		return CreateBookingRequest{}, err
	}
	return request, nil
}
