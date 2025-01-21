package professional

type UpdateAvailabilityRequest struct {
	Items []UpdateAvailabilityItemRequest `json:"items" validate:"required,max=50"`
}

type UpdateAvailabilityItemRequest struct {
	Date string `json:"date" validate:"required,max=50"`
	From string `json:"from" validate:"required,max=50"`
	To   string `json:"to" validate:"required,max=50"`
}
