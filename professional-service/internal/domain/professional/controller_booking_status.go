package professional

import (
	"errors"
	"net/http"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/logger"
	"github.com/jackc/pgx/v5"
)

func (c *Controller) GetBookingStatus(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	ctx := r.Context()
	userId, err := c.userService.GetAuthenticatedUserId(ctx, auth)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	bookingId, err := utils.URLParamInt64(r, "id")
	if err != nil {
		utils.WriteGeneralError(w, errors.New("id is in wrong format"))
		return
	}

	response, err := c.service.GetBookingStatus(ctx, bookingId, userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			utils.WriteJSON(w, http.StatusNotFound, "")
		} else {
			logger.Error("GetBookingStatus", err)
			utils.WriteGeneralError(w, utils.ErrUnknown)
		}
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

type StatusResponse struct {
	Status string `json:"status"`
}
