package booking

import "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/user"

type BookingController struct {
	userService user.Service
	service     *Service
}

func NewController(userService user.Service, service *Service) *BookingController {
	return &BookingController{
		userService: userService,
		service:     service,
	}
}
