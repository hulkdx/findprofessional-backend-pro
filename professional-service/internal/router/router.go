package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
)

func Handler(controller *professional.Controller) http.Handler {
	router := chi.NewRouter()
	router.Get("/professionals", controller.FindAll)
	router.Get("/professional/{id}", controller.Find)
	router.Post("/professional/{id}", controller.Update)
	return router
}
