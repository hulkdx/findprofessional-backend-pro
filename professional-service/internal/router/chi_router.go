package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/controller"
)

func ChiRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/professionals", controller.GetAllProfessionals)
	return router
}
