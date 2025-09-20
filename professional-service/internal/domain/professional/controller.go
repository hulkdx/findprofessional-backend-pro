package professional

import (
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/user"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
)

type Controller struct {
	service      Service
	userService  user.Service
	timeProvider utils.TimeProvider
}

func NewController(service Service, userService user.Service, timeProvider utils.TimeProvider) *Controller {
	return &Controller{
		service:      service,
		userService:  userService,
		timeProvider: timeProvider,
	}
}
