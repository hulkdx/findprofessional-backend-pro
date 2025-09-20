package booking

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/logger"
)

func (c *BookingController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth := r.Header.Get("Authorization")
	userId, err := c.userService.GetAuthenticatedUserId(ctx, auth)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	proID := chi.URLParam(r, "id")
	createBookingRequest, err := parseCreateRequest(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	booking, err := c.service.Create(ctx, userId, proID, createBookingRequest, auth)
	if err != nil {
		switch err {
		case utils.ErrAmountInCentsMismatch, utils.ErrCurrencyMismatch, utils.ErrValidationDatabase:
			utils.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			utils.WriteGeneralError(w, utils.ErrUnknown)
		}
		return
	}
	logger.Debug("booking created:", booking)
	utils.WriteJSON(w, http.StatusOK, booking)
}

func parseCreateRequest(r *http.Request) (booking_model.CreateBookingRequest, error) {
	request := booking_model.CreateBookingRequest{}
	if err := utils.ReadJSON(r, &request); err != nil {
		return booking_model.CreateBookingRequest{}, err
	}
	if err := utils.Validate(request); err != nil {
		return booking_model.CreateBookingRequest{}, err
	}
	return request, nil
}
