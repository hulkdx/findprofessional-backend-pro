package professional

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

var ErrUnknown = errors.New("unknown")

func (c *Controller) Find(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	response, err := c.service.FindById(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusNotFound, "")
		} else {
			utils.WriteGeneralError(w, ErrUnknown)
		}
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}
