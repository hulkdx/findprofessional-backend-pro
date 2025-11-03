package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/data/bookingrepo"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/data/professionalrepo"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking"
	_ "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/booking/model"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/payment"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/user"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewHandler(database *pgxpool.Pool) http.Handler {
	timeProvider := &utils.RealTimeProvider{}
	userService := user.NewService()
	proController := professional.NewController(
		professional.NewService(
			professionalrepo.NewRepository(database, timeProvider),
		),
		userService,
		timeProvider,
	)
	bookingController := booking.NewController(
		userService,
		booking.NewService(
			bookingrepo.NewRepository(database, timeProvider),
			payment.NewService(),
			timeProvider,
		),
	)
	return Handler(proController, bookingController)
}

func Handler(proController *professional.Controller, bookingController *booking.Controller) http.Handler {
	router := chi.NewRouter()

	router.Use(ContentTypeJsonMiddleware)

	normalUser(router, proController)
	proUser(router, proController)

	normalUserBooking(router, bookingController)

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
	// update availability of current pro user
	router.Post("/professional/availability", controller.UpdateAvailability)
}

func normalUserBooking(router *chi.Mux, controller *booking.Controller) {
	//
	// Create a booking for a professional using stripe payment intent
	// ---
	// Request: booking_model.CreateBookingRequest
	// Response: booking_model.CreateBookingResponse
	//
	router.Post("/professional/{id}/booking", controller.Create)
}
