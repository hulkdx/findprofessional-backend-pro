package professional

import (
	"net/http"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

func (c *Controller) GetAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth := r.Header.Get("Authorization")

	userId, err := c.userService.GetAuthenticatedUserId(ctx, auth)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	response, err := c.service.GetAvailability(ctx, userId)
	if err != nil {
		utils.WriteGeneralError(w, utils.ErrUnknown)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}
