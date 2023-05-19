package professional

import (
	"net/http"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

func (c *Controller) FindAll(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if !c.userService.IsAuthenticated(r.Context(), auth) {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorised")
		return
	}
	response, err := c.service.FindAll(r.Context())
	if err != nil {
		utils.WriteGeneralError(w, utils.ErrUnknown)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}
