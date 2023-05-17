package professional

import (
	"database/sql"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/user"
)

type Controller struct {
	service     Service
	userService user.Service
}

func NewController(db *sql.DB) *Controller {
	return &Controller{
		service:     NewService(NewRepository(db)),
		userService: user.NewService(),
	}
}
