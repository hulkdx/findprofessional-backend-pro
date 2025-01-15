package professional

import (
	"net/http"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

func (c *Controller) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth := r.Header.Get("Authorization")

	userId, err := c.userService.GetAuthenticatedUserId(ctx, auth)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	requestBody := UpdateAvailabilityRequest{}
	err = utils.ReadJSON(r, &requestBody)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	err = utils.IsValid(requestBody)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.UpdateAvailability(ctx, userId, requestBody)
	if err != nil {
		utils.WriteGeneralError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "{}")
}
