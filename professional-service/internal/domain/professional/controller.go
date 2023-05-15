package professional

import (
	"database/sql"
)

type Controller struct {
	service Service
}

func NewController(db *sql.DB) *Controller {
	return &Controller{
		service: NewService(NewRepository(db)),
	}
}
