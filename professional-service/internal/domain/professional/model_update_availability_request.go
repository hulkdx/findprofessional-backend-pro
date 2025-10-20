package professional

import "time"

type UpdateAvailabilityRequest struct {
	Items []UpdateAvailabilityItemRequest `json:"items" validate:"required,max=50"`
}

type UpdateAvailabilityItemRequest struct {
	From time.Time `json:"from" validate:"required,datetime"`
	To   time.Time `json:"to" validate:"required,datetime"`
}
