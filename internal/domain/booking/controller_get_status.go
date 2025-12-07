package booking

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

func (c *Controller) GetStatus(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if !c.userService.IsAuthenticated(r.Context(), auth) {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	bookingId := chi.URLParam(r, "id")
	status, err := c.service.GetStatus(r.Context(), bookingId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusNotFound, "")
		} else {
			utils.WriteGeneralError(w, utils.ErrUnknown)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, status)
}
