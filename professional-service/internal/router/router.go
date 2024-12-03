package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
)

func Handler(controller *professional.Controller) http.Handler {
	router := chi.NewRouter()

	router.Use(ContentTypeJsonMiddleware)

	router.Get("/professional", controller.FindAll)
	router.Put("/professional", controller.Create)
	router.Get("/professional/{id}", controller.Find)
	router.Post("/professional", controller.Update)

	router.Get("/professional/{id}/review", controller.FindAllReview)

	return router
}
