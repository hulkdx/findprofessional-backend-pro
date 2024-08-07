package professional

import (
	"net/http"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/lib/pq"
)

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	request := CreateRequest{}
	if err := utils.ReadJSON(r, &request); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := utils.IsValid(request); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.service.Create(r.Context(), request); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			utils.WriteError(w, http.StatusConflict, "")
		} else {
			utils.WriteGeneralError(w, utils.ErrUnknown)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
