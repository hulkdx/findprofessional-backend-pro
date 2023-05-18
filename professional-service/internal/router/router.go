package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
)

func Handler(controller *professional.Controller) http.Handler {
	router := chi.NewRouter()
	router.Get("/professionals", controller.FindAll)
	return router
}
