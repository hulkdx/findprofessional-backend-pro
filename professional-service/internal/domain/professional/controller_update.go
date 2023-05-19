package professional

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if !c.userService.IsAuthenticated(r.Context(), auth) {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	id := chi.URLParam(r, "id")

	updateRequest := UpdateRequest{}
	err := utils.ReadJSON(r, &updateRequest)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = c.service.Update(r.Context(), id, updateRequest)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusNotFound, "")
		} else {
			utils.WriteGeneralError(w, utils.ErrUnknown)
		}
		return
	}
	idInt, _ := strconv.Atoi(id)
	response := Professional{
		ID:        idInt,
		Email:     updateRequest.Email,
		UpdatedAt: nil,
		CreatedAt: nil,
	}
	utils.WriteJSON(w, http.StatusOK, response)
}
