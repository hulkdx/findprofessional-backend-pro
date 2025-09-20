package professional

import (
	"net/http"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request, err := parseCreateRequest(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.Create(ctx, request)
	if err != nil {
		switch err {
		case utils.ErrDuplicate:
			utils.WriteError(w, http.StatusConflict, "Duplicate user")
		case utils.ErrNotFoundUser:
			utils.WriteError(w, http.StatusNotFound, "Not found user")
		default:
			utils.WriteGeneralError(w, utils.ErrUnknown)
		}
		return
	}

	response, err := c.userService.Login(ctx, request.Email, request.Password)
	if err != nil {
		utils.WriteGeneralError(w, utils.ErrUnknown)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func parseCreateRequest(r *http.Request) (CreateRequest, error) {
	request := CreateRequest{}
	if err := utils.ReadJSON(r, &request); err != nil {
		return CreateRequest{}, err
	}
	if err := utils.Validate(&request); err != nil {
		return CreateRequest{}, err
	}
	return request, nil
}
