package booking

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/logger"
)

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth := r.Header.Get("Authorization")
	userId, err := c.userService.GetAuthenticatedUserId(ctx, auth)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	proIdStr := chi.URLParam(r, "id")
	proId, err := strconv.ParseInt(proIdStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	request, err := bookingmodel.ParseCreateRequest(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(request.Availabilities) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "no availabilities")
		return
	}
	if len(request.Availabilities) > 48 {
		utils.WriteError(w, http.StatusBadRequest, "invalid availability size > 48")
		return
	}

	params := CreateParams{
		ProId:          proId,
		UserId:         userId,
		AmountInCents:  request.AmountInCents,
		Currency:       request.Currency,
		IdempotencyKey: request.IdempotencyKey,
		Auth:           auth,
		Availabilities: request.Availabilities,
	}
	response, err := c.service.Create(ctx, &params)
	if err != nil {
		var safe *utils.SafeHttpError
		if errors.As(err, &safe) {
			utils.WriteError(w, http.StatusBadRequest, err.Error())
		} else {
			utils.WriteGeneralError(w, utils.ErrUnknown)
		}
		return
	}
	logger.Debug("response created:", response)
	utils.WriteJSON(w, http.StatusOK, response)
}
