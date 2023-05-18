package professional

import (
	"database/sql"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/user"
)

type Controller struct {
	service     Service
	userService user.Service
}

func NewController(service Service, userService user.Service) *Controller {
	return &Controller{
		service:     service,
		userService: userService,
	}
}

func NewControllerFromDB(db *sql.DB) *Controller {
	return &Controller{
		service:     NewService(NewRepository(db)),
		userService: user.NewService(),
	}
}
