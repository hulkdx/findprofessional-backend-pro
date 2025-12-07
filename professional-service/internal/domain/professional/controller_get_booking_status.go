package professional

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/jackc/pgx/v5"
)

func (c *Controller) GetBookingStatus(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	ctx := r.Context()
	if !c.userService.IsAuthenticated(ctx, auth) {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	bookingIdStr := chi.URLParam(r, "id")
	bookingId, err := strconv.ParseInt(bookingIdStr, 10, 64)
	if err != nil {
		utils.WriteGeneralError(w, errors.New("id is in wrong format"))
		return
	}

	response, err := c.service.GetBookingStatus(ctx, bookingId)
	if err != nil {
		if err == pgx.ErrNoRows {
			utils.WriteJSON(w, http.StatusNotFound, "")
		} else {
			utils.WriteGeneralError(w, utils.ErrUnknown)
		}
		return
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

type StatusResponse struct {
	Status string `json:"status"`
}
