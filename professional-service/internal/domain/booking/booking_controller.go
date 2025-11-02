package booking

import "github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/user"

type Controller struct {
	userService user.Service
	service     *Service
}

func NewController(userService user.Service, service *Service) *Controller {
	return &Controller{
		userService: userService,
		service:     service,
	}
}
