package professional

import (
	"database/sql"

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

func NewControllerFromDB(db *sql.DB) *Controller {
	timeProvider := &utils.RealTimeProvider{}
	return &Controller{
		service:      NewService(NewRepository(db, timeProvider)),
		userService:  user.NewService(),
		timeProvider: timeProvider,
	}
}
