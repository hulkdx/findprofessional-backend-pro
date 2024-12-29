package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
)

func Handler(controller *professional.Controller) http.Handler {
	router := chi.NewRouter()

	router.Use(ContentTypeJsonMiddleware)

	normalUser(router, controller)
	proUser(router, controller)

	return router
}

func normalUser(router *chi.Mux, controller *professional.Controller) {
	router.Get("/professional", controller.FindAll)
	router.Get("/professional/{id}", controller.Find)
	router.Get("/professional/{id}/review", controller.FindAllReview)
}

func proUser(router *chi.Mux, controller *professional.Controller) {
	// register a new pro
	router.Put("/professional", controller.Create)
	// update pro profile information
	router.Post("/professional", controller.Update)
	// get availability information from time schedules of current pro user
	router.Get("/professional/availability", controller.GetAvailability)
}
