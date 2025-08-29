package professional

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth := r.Header.Get("Authorization")

	userId, err := c.userService.GetAuthenticatedUserId(ctx, auth)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	updateRequest := UpdateRequest{}
	err = utils.ReadJSON(r, &updateRequest)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	err = utils.Validate(updateRequest)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.Update(ctx, strconv.FormatInt(userId, 10), updateRequest)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusNotFound, "")
		} else {
			utils.WriteGeneralError(w, utils.ErrUnknown)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, "")
}
