package professional

import (
	"net/http"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/logger"
)

func (c *Controller) GetBookingUser(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	ctx := r.Context()
	userId, err := c.userService.GetAuthenticatedUserId(ctx, auth)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorised")
		return
	}
	response, err := c.service.GetBookings(ctx, userId, UserTypeNormal)
	if err != nil {
		logger.Error("GetBookingUser", err)
		utils.WriteGeneralError(w, utils.ErrUnknown)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}
