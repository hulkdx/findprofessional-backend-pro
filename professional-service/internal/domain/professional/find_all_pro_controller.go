package professional

import (
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/logger"
	"net/http"
)

func (c *Controller) FindAllProfessional(w http.ResponseWriter, r *http.Request) {
	// TODO: authenticate
	response, err := c.service.FindAllProfessional()
	if err != nil {
		// TODO:
		w.Write([]byte(err.Error()))
		return
	}
	err = utils.WriteJSON(w, http.StatusOK, response)
	if err != nil {
		logger.Error("WriteJSON error: ", err)
	}
}
