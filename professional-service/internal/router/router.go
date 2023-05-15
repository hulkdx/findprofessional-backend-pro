package router

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
)

func Handler(db *sql.DB) http.Handler {
	controller := professional.NewController(db)

	router := chi.NewRouter()
	router.Get("/professionals", controller.GetAllProfessionals)
	return router
}
