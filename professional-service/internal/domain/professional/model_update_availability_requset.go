package professional

type UpdateAvailabilityRequest struct {
	Items []UpdateAvailabilityItemRequest `json:"items" validate:"max=50"`
}

type UpdateAvailabilityItemRequest struct {
	Date string `json:"date" validate:"max=50"`
	From string `json:"from" validate:"max=50"`
	To   string `json:"to" validate:"max=50"`
}