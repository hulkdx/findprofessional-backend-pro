package professional

import (
	"errors"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"net/http"
)

func (c *Controller) FindAllProfessional(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if !c.userService.IsAuthenticated(r.Context(), auth) {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorised")
		return
	}
	response, err := c.service.FindAllProfessional(r.Context())
	if err != nil {
		utils.WriteGeneralError(w, errors.New("invalid data"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}
