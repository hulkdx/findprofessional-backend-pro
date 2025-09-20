package booking

import "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/user"

type BookingController struct {
	userService user.Service
	service     *BookingService
}

func NewController(userService user.Service, service *BookingService) *BookingController {
	return &BookingController{
		userService: userService,
		service:     service,
	}
}
