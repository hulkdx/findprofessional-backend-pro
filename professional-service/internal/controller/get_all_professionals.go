package controller

import "net/http"

func GetAllProfessionals(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("[]"))
}
