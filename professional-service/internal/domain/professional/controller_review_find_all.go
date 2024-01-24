package professional

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

func (c *Controller) FindAllReview(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if !c.userService.IsAuthenticated(r.Context(), auth) {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	professionalIdStr := chi.URLParam(r, "id")
	professionalId, err := strconv.ParseInt(professionalIdStr, 10, 64)
	if err != nil {
		utils.WriteGeneralError(w, errors.New("id is in wrong format"))
		return
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		utils.WriteGeneralError(w, errors.New("page is in wrong format"))
		return
	}

	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		utils.WriteGeneralError(w, errors.New("pageSize is in wrong format"))
		return
	}

	response, err := c.service.FindAllReview(r.Context(), professionalId, page, pageSize)
	if err != nil {
		utils.WriteGeneralError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}
