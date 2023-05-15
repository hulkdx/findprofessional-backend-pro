package professional

import (
	"net/http"
)

func (c *Controller) GetAllProfessionals(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("[]"))
}
