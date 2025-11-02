package booking

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	booking_model "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/logger"
)

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth := r.Header.Get("Authorization")
	userId, err := c.userService.GetAuthenticatedUserId(ctx, auth)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	proID := chi.URLParam(r, "id")
	createBookingRequest, err := booking_model.ParseCreateRequest(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	booking, err := c.service.Create(ctx, userId, proID, createBookingRequest, auth)
	if err != nil {
		var safe *utils.SafeHttpError
		if errors.As(err, &safe) {
			utils.WriteError(w, http.StatusBadRequest, err.Error())
		} else {
			utils.WriteGeneralError(w, utils.ErrUnknown)
		}
		return
	}
	logger.Debug("booking created:", booking)
	utils.WriteJSON(w, http.StatusOK, booking)
}
